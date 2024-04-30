package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/logger"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var consulAddress = "test.site.com"
var consulClient *api.Client
var serviceId string

func NewGrpcServer() *grpc.Server {
	// TODO 定义grpc监控指标
	// grpc_prometheus.GrpcServerMonitorMetricsRegist()
	// grpc_prometheus.GetGrpcServerMetrics().EnableHandlingTimeSummary()

	s := grpc.NewServer(
		grpc.MaxConcurrentStreams(100),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),

		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle: time.Duration(30) * time.Second,
				Time:              time.Duration(15) * time.Second,
				Timeout:           time.Duration(3) * time.Second,
			},
		),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				//grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpcRecovery.UnaryServerInterceptor(),
				grpcOpentracing.UnaryServerInterceptor(),
				//grpc_prometheus.GetGrpcServerMetrics().UnaryServerInterceptor(),
			),
		),
	)

	return s
}

func GetGrpcConn(ctx context.Context, name string) (*IdleConn, context.Context) {
	newCtx := setTraceCtx(ctx)
	conn, err := grpcConnPool[name].Get()
	if err == nil {
		defer grpcConnPool[name].Put(conn)
	}

	return conn, newCtx
}

func PutGrpcConn(name string, conn *IdleConn) {
	grpcConnPool[name].RetrieveConcurrentStream(conn)
}

// setTraceCtx 设置trace_id, span_id等信息
func setTraceCtx(ctx context.Context) context.Context {
	var ctxMD = ctx
	//
	//requestId := ctx.Value("x-request-id")
	//
	//traceId := ctx.Value("X-B3-TraceId")
	//spanId := ctx.Value("X-B3-SpanId")
	////parentSpanId := ctx.GetHeader("x-b3-parentspanid")
	//sampled := ctx.Value("X-B3-Sampled")
	//
	//if sampled == "1" && traceId != "" {
	//	number.GenerateTraceId()
	//	md := metadata.Pairs("X-B3-Sampled", sampled.(string),
	//		"x-request-id", requestId.(string),
	//		"X-B3-TraceId", traceId.(string),
	//		"X-B3-SpanId", number.GenerateSpanId(8),
	//		"X-B3-ParentSpanId", spanId.(string),
	//		"x-b3-sampled", sampled.(string),
	//		"x-b3-flags", ctx.Value("x-b3-flags").(string),
	//		"x-ot-span-context", ctx.Value("x-ot-span-context").(string),
	//		//"X-B3-ProjectId", ctx.GetHeader("X-B3-ProjectId"),
	//	)
	//	ctxMD = metadata.NewOutgoingContext(context.Background(), md)
	//} else {
	//	ctxMD = context.Background()
	//}

	return ctxMD
}

func panicLog(err error) {
	CoreCtx.Logger.Error(err.Error(), zap.String("type", "grpc"))
}

func RunServer(s *grpc.Server, addr string) {
	// register reflection
	registerGrpcReflection(s)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("grpc", zap.NamedError("grpc", errors.New("grpc error"+err.Error())))
	}

	if err != nil {
		os.Exit(1)
	}
	fmt.Println("local grpc start on " + addr)

	go s.Serve(l)
	port := 0
	address := os.Getenv("HOST_IP")
	hostPort := os.Getenv("HOST_PORT")
	if hostPort != "" {
		port, _ = strconv.Atoi(hostPort)
	}
	// 将当前grpc服务注册到consul
	err = RegisterGRPCServiceConsul(address, port, []string{GetProjectName()})
	if err != nil {
		logger.Error("grpc", zap.NamedError("grpc", errors.New("grpc register error "+err.Error())))
	}
	gracefulGrpcShutdown(s)
}

func getLocalIp() string {
	conn, err := net.Dial("udp", "183.60.83.19:53")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return strings.Split(localAddr.String(), ":")[0]

}

// registerGrpcReflection 用户本地调试
func registerGrpcReflection(s *grpc.Server) {
	if DevEnv == EnvLocal {
		reflection.Register(s)
	}
}

func gracefulGrpcShutdown(s *grpc.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, os.Interrupt)
	getSignal := <-ch
	if getSignal == syscall.SIGTERM || getSignal == syscall.SIGQUIT || getSignal == syscall.SIGINT {

		logger.Debug("grpc", zap.NamedError("shutdown", errors.New("grpc server shutdown, handle signal "+fmt.Sprint(getSignal))))
		time.Sleep(3 * time.Second)
	}
	s.GracefulStop()
	s.Stop()
	if consulClient != nil {
		consulClient.Agent().ServiceDeregister(serviceId)
	}

	os.Exit(0)
}

// RegisterGRPCServiceConsul 注册服务到consul
func RegisterGRPCServiceConsul(address string, port int, tags []string) error {
	serviceId = uuid.NewV4().String()
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", consulAddress, 8500)
	consulClient, err = api.NewClient(cfg)
	if err != nil {
		logger.Error("grpc", zap.NamedError("grpc", errors.New("grpc register to consul error "+err.Error())))
	}

	// 生成健康检查对象
	check := &api.AgentServiceCheck{
		GRPC: fmt.Sprintf("%s:%d", address, port), // 服务的运行地址,ip不可以是127.0.0.1

		Timeout:                        "3s",  // 超过此时间说明服务状态不健康
		Interval:                       "5s",  // 每5s检查一次
		DeregisterCriticalServiceAfter: "30s", // 失败多久后注销服务
	}
	// 生成注册对象
	serverName := fmt.Sprintf("%s-%s", DevEnv, GetProjectName())
	registration := &api.AgentServiceRegistration{
		Name:    serverName,
		ID:      serviceId,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check:   check,
	}

	// 注册服务
	err = consulClient.Agent().ServiceRegister(registration)

	return err

}
