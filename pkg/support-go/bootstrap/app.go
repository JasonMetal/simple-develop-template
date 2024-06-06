package bootstrap

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/JasonMetal/simple-develop-template/pkg/support-go/helper/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	OnlineServiceHostPort = ":50051"
	HttpTimeoutHandler    = 10 * time.Second //TimeoutHandler默认的http超时时间，http响应时间超过后直接返回client端503(不同项目可根据接口最大超时时间调整)
	EnvLocal              = "local"
	EnvTest               = "test"
	EnvBVT                = "bvt"
	EnvProduct            = "product"
	EnvBenchmark          = "benchmark"
)

type coreCtx struct {
	//ConsoleLogger *zap.Logger
	Logger    *zap.Logger
	UDPLogger *zap.Logger
	//TracingLogKafkaCollect collect.Collector
	//NoCallerLogger *zap.Logger
}

var (
	CoreCtx     coreCtx
	err         error
	DevEnv      string
	TestConfig  string
	ProjectName string
)

//var grpcConnPool map[string]http2.Pool

func Init() {
	initEnv()

	InitLogger()

	InitMysql()

	InitRedis()

	InitGrpc()

	// InitSts()

	// InitSms()
}

func SetProjectName(name string) {
	ProjectName = name
}

func GetProjectName() string {
	return ProjectName
}

func initEnv() {
	// 优先获取系统环境变量
	runEnv := os.Getenv("RUN_ENV")
	if runEnv != "" {
		DevEnv = runEnv
	} else {
		flag.Usage = Usage
		flag.StringVar(&DevEnv, "e", EnvLocal, "Specify env")
		flag.Parse()
	}
}

func InitWeb(funs []gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(SetLogger())
	r.Use(LogRecovery())
	r.Use(gin.Recovery())
	r.Use(ControlCors())

	for _, v := range funs {
		r.Use(v)
	}

	if DevEnv == EnvLocal || DevEnv == EnvBenchmark {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGKILL)
		go func() {
			<-c
			os.Exit(0)
		}()
	} else {
		logger.RedirectLog()
	}

	return r
}

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-e=local]")
	flag.PrintDefaults()
	os.Exit(0)
}

func RunWeb(r *gin.Engine, addr string) {
	if DevEnv == EnvLocal || DevEnv == EnvBenchmark {
		fmt.Println("local http start on " + addr)
		//router.Run(addr)
	} else {
		fmt.Printf("env:%s; http server listen on %s\n", DevEnv, OnlineServiceHostPort)
		addr = OnlineServiceHostPort
		//router.Run(core.ONLINE_SERVICE_HOST_PORT)
	}
	s := &http.Server{
		Addr: addr,
		Handler: http.TimeoutHandler(
			r,
			HttpTimeoutHandler,
			"server has gone away",
		),
		ReadTimeout:  HttpTimeoutHandler,
		WriteTimeout: HttpTimeoutHandler,
		IdleTimeout:  1 * time.Minute,
		//MaxHeaderBytes: 1 << 20,
	}

	go s.ListenAndServe()

	gracefulShutdown(s)
}

// gracefulShutdown 优雅退出
func gracefulShutdown(server *http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, os.Interrupt)
	<-ch
	//core.DebugLog("http shutdown")

	cxt, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	server.Shutdown(cxt)
	os.Exit(0)

}

// ControlCors 设置CORS
func ControlCors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")

		if strings.Contains(origin, "sitename.cloud") || strings.Contains(origin, "localhost") {
			context.Header("Access-Control-Allow-Origin", origin)
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			context.Header("Access-Control-Allow-Headers", "*")
			context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, admintoken")
			context.Header("Access-Control-Allow-Credentials", "true")

			if method == "OPTIONS" {
				context.AbortWithStatus(http.StatusNoContent)
				return
			}
		}
		context.Next()
	}
}

func CheckError(err error, errType string) error {
	if err != nil {

		logger.Error(err.Error(), zap.String("type", errType))

		return err
	}

	return nil
}

func ProjectPath() (path string) {
	// default linux/mac os
	var (
		sp = "/"
		ss []string
	)
	if runtime.GOOS == "windows" {
		sp = "\\"
	}

	// GOMOD
	// in go source code:
	// // Check for use of modules by 'go env GOMOD',
	// // which reports a go.mod file path if modules are enabled.
	// stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	// gomod := string(bytes.TrimSpace(stdout))
	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	path = string(bytes.TrimSpace(stdout))
	if path == "/dev/null" {
		return ""
	}
	if path != "" {
		ss = strings.Split(path, sp)
		ss = ss[:len(ss)-1]
		path = strings.Join(ss, sp) + sp
		return
	}

	// GOPATH
	fileDir, _ := os.Getwd()
	path = os.Getenv("GOPATH") // < go 1.17 use
	ss = strings.Split(fileDir, path)
	if path != "" {
		ss2 := strings.Split(ss[1], sp)
		path += sp
		for i := 1; i < len(ss2); i++ {
			path += ss2[i] + sp
			return path
		}
	}
	return
}
