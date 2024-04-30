package bootstrap

import (
	"fmt"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/config"
	"github.com/BioforestChain/simple-develop-template/pkg/support-go/helper/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type loggerConfig struct {
	Level      string
	LogType    string
	LogPath    string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

func InitLogger() {
	path := fmt.Sprintf("%smanifest/config/%s/logger.yml", ProjectPath(), DevEnv)

	cfg, err := config.GetConfig(path)
	if err != nil {
		panic("init logger err" + err.Error())
	}

	filename, _ := cfg.String("logger.filename")
	maxSize, _ := cfg.Int("logger.maxSize")
	maxBackup, _ := cfg.Int("logger.maxBackup")
	maxAge, _ := cfg.Int("logger.maxAge")
	compress, _ := cfg.Bool("logger.compress")
	logType, _ := cfg.String("logger.logType")
	level, _ := cfg.String("logger.level")
	logPath, _ := cfg.String("logger.logPath")
	logcfg := &loggerConfig{
		Level:      level,
		LogType:    logType,
		LogPath:    logPath,
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
	// 获取日志写入介质
	writeSyncer := getLogWriter(logcfg)
	// 设置日志等级，具体请见 config/logger.yml 文件
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误，日志级别设置有误。请修改 config/logger.yml 文件中的 logger.level 配置项")
	}
	//初始化错误等级
	_ = logLevel.Set(level)
	// 初始化 core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 初始化 Logger
	Logger := zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号，内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层，调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error 时才会显示 stacktrace
	)

	// 将自定义的 logger 替换为全局的 logger
	logger.SetLogger(Logger)

	// zap.L().Fatal() 调用时，就会使用我们自定的 Logger
	zap.ReplaceGlobals(Logger)

	return
}

// getEncoder 设置日志存储格式
func getEncoder() zapcore.Encoder {

	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用，如 paginator/paginator.go:148
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,     // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,   // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     customTimeEncoder,             // 时间格式，我们自定义为 2006-01-02 15:04:05
		EncodeDuration: zapcore.MillisDurationEncoder, // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,    // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	// 本地环境配置
	if DevEnv == EnvLocal {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// 本地设置内置的 Console 解码器（支持 stacktrace 换行）
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 线上环境使用 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// customTimeEncoder 自定义友好的时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// getLogWriter 日志记录介质
func getLogWriter(logcfg *loggerConfig) zapcore.WriteSyncer {
	filename := logcfg.Filename
	// 按照日期记录日志文件
	if logcfg.LogType == "daily" {
		filename = logcfg.LogPath + GetProjectName() + "-" + time.Now().Format("2006-01-02.log")
	}

	// 滚动日志，详见 config/logger.yml
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    logcfg.MaxSize,
		MaxBackups: logcfg.MaxBackups,
		MaxAge:     logcfg.MaxAge,
		Compress:   logcfg.Compress,
	}
	// 配置输出介质
	if DevEnv == EnvLocal {
		// 本地开发终端打印和记录文件
		//return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
		return zapcore.AddSync(lumberJackLogger)

	} else {
		// 生产环境只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}
}

// SetLogger 接收gin框架默认的日志
func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		// 过滤健康检测
		if c.Request.UserAgent() != "clb-healthcheck" && path != "/ping" {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("cost", cost),
			)
		}
	}
}

func LogRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				c.JSON(http.StatusInternalServerError, map[string]any{
					"code":    1000,
					"data":    struct{}{},
					"message": "api panic",
				})
				c.AbortWithStatus(http.StatusInternalServerError)

			}
		}()
		c.Next()
	}
}
