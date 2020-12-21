package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
)

// LogPath log path
var LogPath string = "/tmp/aphrodite-web/access.log"

// ip , time, method, path, proto , statusCode, Latency, UserAgent, ErrorMessage
var format = "%s -- [%s] \"%s %s %s %d %s \"%s\" %s\"\n"

// LogWriter log writer
func LogWriter() io.Writer {
	writer, err := rotatelogs.New(
		LogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(LogPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationSize(1024*1024*5),
	)
	if err != nil {
		zlog.Fatal("Panic: Init access.log failure", zap.Error(err))
	}
	return writer
}

// LogInterceptor log interceptor
func LogInterceptor() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf(format,
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
