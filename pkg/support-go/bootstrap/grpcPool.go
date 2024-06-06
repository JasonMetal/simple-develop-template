package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/config"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	yCfg "github.com/olebedev/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrClosed            = errors.New("pool had closed")
	ErrLimit             = errors.New("pool is reached max")
	MaxConcurrentStreams = 100
	MaxGetConnTimes      = 200
)

// Pool 定义pool接口
type Pool interface {
	Get() (*IdleConn, error)

	Put(*IdleConn) error

	RetrieveConcurrentStream(*IdleConn)

	RetrieveStream(*IdleConn)

	Close(*grpc.ClientConn) error

	Release()

	Len() int
}

type PoolConfig struct {
	InitialCap  int
	MaxCap      int
	Factory     func(string, time.Duration) (*grpc.ClientConn, error)
	Close       func(*grpc.ClientConn) error
	IdleTimeout time.Duration
	Endpoint    string
	Timeout     time.Duration
}

type ChannelPool struct {
	mu          sync.Mutex
	conns       chan *IdleConn
	factory     func(string, time.Duration) (*grpc.ClientConn, error)
	close       func(*grpc.ClientConn) error
	idleTimeout time.Duration
	initialCap  int
	maxCap      int
	endpoint    string
	timeout     time.Duration
}

type IdleConn struct {
	Conn          *grpc.ClientConn
	t             time.Time
	currentStream int32
}

type GrpcEndpoint struct {
	Endpoint    string
	InitialCap  int
	MaxCap      int
	IdleTimeout int
}

var grpcConnPool map[string]Pool

// InitGrpc 初始化连接池
func InitGrpc() {
	var endpoints = make(map[string]GrpcEndpoint)

	path := fmt.Sprintf("%smanifest/config/%s/grpc.yml", ProjectPath(), DevEnv)
	cfg, err := config.GetConfig(path)
	if err != nil {
		panic("init grpc err")
	}
	servers, err := cfg.List("server")
	for _, v := range servers {
		name, _ := yCfg.Get(v, "name")
		endpoint, _ := yCfg.Get(v, "endpoint")
		initialCap, _ := yCfg.Get(v, "initialCap")
		maxCap, _ := yCfg.Get(v, "maxCap")
		idleTimeout, _ := yCfg.Get(v, "idleTimeout")

		endpoints[(name.(string))] = GrpcEndpoint{
			endpoint.(string),
			initialCap.(int),
			maxCap.(int),
			idleTimeout.(int),
		}
	}

	grpcConnPool = NewPool(endpoints)

	endpoints = nil
}

// NewChannelPool 初始化连接池
func NewChannelPool(poolConfig *PoolConfig) (Pool, error) {
	if poolConfig.InitialCap < 0 || poolConfig.MaxCap <= 0 || poolConfig.InitialCap > poolConfig.MaxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &ChannelPool{
		conns:       make(chan *IdleConn, poolConfig.MaxCap),
		factory:     poolConfig.Factory,
		close:       poolConfig.Close,
		idleTimeout: poolConfig.IdleTimeout,
		initialCap:  poolConfig.InitialCap,
		maxCap:      poolConfig.MaxCap,
		endpoint:    poolConfig.Endpoint,
		timeout:     poolConfig.Timeout,
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory(poolConfig.Endpoint, poolConfig.Timeout)
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &IdleConn{Conn: conn, t: time.Now()}
	}

	return c, nil
}

// getConns 获取连接池
func (c *ChannelPool) getConns() chan *IdleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()

	return conns
}

// Get 从连接池获取连接
func (c *ChannelPool) Get() (*IdleConn, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	size := len(conns)
	times := 0
	i := 0
	//var current int32 = 0

	for {
		select {
		case wrapConn := <-conns:
			times++
			if wrapConn == nil {
				return nil, ErrClosed
			}

			if size > c.initialCap {
				if timeout := c.idleTimeout; timeout > 0 {
					if wrapConn.t.Add(timeout).Before(time.Now()) {
						i++
						c.Close(wrapConn.Conn)
						continue
					}
				}
			}

			grpcState := wrapConn.Conn.GetState()
			if grpcState == connectivity.TransientFailure || grpcState == connectivity.Shutdown {
				i++
				log.Printf("grpc recycle state:%d, current size:%d, info:%s", grpcState, len(c.conns), grpcState.String())
				c.Close(wrapConn.Conn)
				continue
			}

			if atomic.LoadInt32(&wrapConn.currentStream) < int32(MaxConcurrentStreams) {
				atomic.AddInt32(&wrapConn.currentStream, 1)

				return wrapConn, nil
			} else {
				c.conns <- wrapConn
				i++

				if i > size && c.maxCap > c.Len() {
					conn, err := c.factory(c.endpoint, c.timeout)
					if err != nil {
						return nil, err
					}
					idleConn := &IdleConn{Conn: conn, t: time.Now(), currentStream: 1}

					return idleConn, nil
				} else if times > MaxGetConnTimes {
					return nil, ErrLimit
				}
			}

		default:
			if c.maxCap > c.Len() {
				conn, err := c.factory(c.endpoint, c.timeout)
				if err != nil {
					return nil, err
				}
				idleConn := &IdleConn{Conn: conn, t: time.Now(), currentStream: 1}

				return idleConn, nil
			}
			return nil, ErrLimit
		}
	}
}

// Put 连接放回连接池
func (c *ChannelPool) Put(ic *IdleConn) error {
	if ic.Conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	grpcState := ic.Conn.GetState()
	if grpcState == connectivity.TransientFailure || grpcState == connectivity.Shutdown {
		c.Close(ic.Conn)
		log.Printf("grpc recycle state:%d, current size:%d", ic.Conn.GetState(), len(c.conns))

		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		return c.Close(ic.Conn)
	}

	ic.t = time.Now()

	select {
	case c.conns <- ic:
		return nil
	default:
		return c.Close(ic.Conn)
	}
}

// RetrieveConcurrentStream 回收steam
func (c *ChannelPool) RetrieveConcurrentStream(ic *IdleConn) {
	atomic.AddInt32(&ic.currentStream, -1)
}

// RetrieveStream 回收steam
func (c *ChannelPool) RetrieveStream(ic *IdleConn) {
	atomic.AddInt32(&ic.currentStream, -1)
	c.Put(ic)
}

// Close 关闭链接
func (c *ChannelPool) Close(conn *grpc.ClientConn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.close(conn)
}

// Release 释放所有的链接
func (c *ChannelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		closeFun(wrapConn.Conn)
	}
}

// Len 连接池中已有的连接
func (c *ChannelPool) Len() int {
	return len(c.getConns())
}

func NewPool(endpoints map[string]GrpcEndpoint) map[string]Pool {
	var connPool = make(map[string]Pool)
	var timeout time.Duration = time.Duration(4000) * time.Millisecond
	for name, endpoint := range endpoints {
		factory := func(endpoint string, timeout time.Duration) (*grpc.ClientConn, error) {
			return newConn(endpoint, timeout)
		}
		closeFunc := func(conn *grpc.ClientConn) error { return conn.Close() }
		poolcfg := &PoolConfig{
			InitialCap:  endpoint.InitialCap,
			MaxCap:      endpoint.MaxCap,
			Factory:     factory,
			Close:       closeFunc,
			IdleTimeout: time.Duration(endpoint.IdleTimeout) * time.Second,
			Endpoint:    endpoint.Endpoint,
			Timeout:     timeout,
		}
		var err error
		connPool[name], err = NewChannelPool(poolcfg)
		// todo log grpc error
		if err != nil {

		}
	}
	endpoints = nil

	return connPool
}

func newConn(endpoint string, timeout time.Duration) (*grpc.ClientConn, error) {
	retryOpts := []grpc_retry.CallOption{
		//grpc_retry.WithCodes(codes.Canceled, codes.DataLoss, codes.Unavailable),
		grpc_retry.WithCodes(codes.DataLoss, codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // returns a DialOption which disables transport security for this ClientConn
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(15) * time.Second,
			Timeout:             time.Duration(3) * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			//grpc_opentracing.UnaryClientInterceptor(),
			TimeoutUnaryClientInterceptor(timeout),
			grpc_retry.UnaryClientInterceptor(retryOpts...),
			RPCConnectMonitor(),
			//grpc_prometheus.GetGrpcClientMetrics().UnaryClientInterceptor(),
		)),
	)

	// todo log errors

	return conn, err
}

// RPCConnectMonitor gRPC.Server 重启, gRPC.Pool连接的时健康检测
func RPCConnectMonitor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		grpcState := cc.GetState()

		if grpcState == connectivity.TransientFailure || grpcState == connectivity.Shutdown {
			cc.ResetConnectBackoff()
			grpcState = cc.GetState()
		}

		if grpcState != connectivity.Ready && grpcState != connectivity.Connecting && grpcState != connectivity.Idle {
			//log.Printf("grpc.Server unavailable, state: %d, method: %s, req: %v", grpcState, method, req)
			return status.Error(codes.Unavailable, "")
		}

		//newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		//defer cancel()

		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

// TimeoutUnaryClientInterceptor returns a new unary client interceptor that sets a timeout on the request context.
func TimeoutUnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		deadline, _ := ctx.Deadline()
		if time.Until(deadline).Microseconds() > 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctxTm, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(ctxTm, method, req, reply, cc, opts...)
	}
}
