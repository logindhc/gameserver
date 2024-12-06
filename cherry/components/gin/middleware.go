// Package cherryGin from https://github.com/gin-contrib/zap/
package cherryGin

import (
	"fmt"
	cherryCrypto "gameserver/cherry/extend/crypto"
	cstring "gameserver/cherry/extend/string"
	csync "gameserver/cherry/extend/sync"
	clog "gameserver/cherry/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func GinDefaultZap() GinHandlerFunc {
	return GinZap(time.RFC3339, true)
}

// GinZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//  1. A time package format string (e.g. time.RFC3339).
//  2. A boolean stating whether to use UTC time zone or local.
func GinZap(timeFormat string, utc bool) GinHandlerFunc {
	return func(c *Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				clog.Error(e)
			}
		} else {
			clog.Debugw(c.FullPath(),
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", path,
				"query", query,
				"ip", c.ClientIP(),
				"user-agent", c.Request.UserAgent(),
				"time", end.Format(timeFormat),
				"latency", latency,
			)
		}
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithZap(stack bool) GinHandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					clog.Warnw(c.Request.URL.Path,
						"error", err,
						"request", string(httpRequest),
					)

					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					clog.Warnw("[Recovery from panic]",
						"time", time.Now(),
						"error", err,
						"request", string(httpRequest),
						"stack", string(debug.Stack()),
					)
				} else {
					clog.Warnw("[Recovery from panic]",
						"time", time.Now(),
						"error", err,
						"request", string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 开始时间
		startTime := time.Now()
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 消耗的时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 用户IP
		clientIP := c.ClientIP()
		// 日志格式
		clog.Debugf("| %3d | %13v | %15s | %s %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func Cors(domain ...string) GinHandlerFunc {
	return func(c *Context) {
		method := c.Request.Method

		if len(domain) > 0 {
			c.Header("Access-Control-Allow-Origin", domain[0])
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func Md5Filter(md5Key string) GinHandlerFunc {
	return func(c *Context) {
		if cstring.IsBlank(md5Key) {
			//不配置就不做验证md5
			c.Next()
			return
		}
		sign := c.GetString("sign", "", true)
		if cstring.IsBlank(sign) {
			clog.Debug("req param sign is nil")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		timeStamp := c.GetInt64("time", 0, true)
		if timeStamp <= 0 {
			clog.Debug("req param time is nil or 0")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		md5 := cherryCrypto.MD5(fmt.Sprintf("%v%s", timeStamp, md5Key))
		if sign != md5 {
			clog.Debugf("md5 check failed, sign = %s, md5 = %s", sign, md5)
			c.AbortWithStatus(http.StatusForbidden)
		}
		c.Next()
	}
}

// MaxConnect limit max connect
func MaxConnect(n int) GinHandlerFunc {
	latch := csync.NewLimit(n)

	return func(c *Context) {
		if latch.TryBorrow() {
			defer func() {
				if err := latch.Return(); err != nil {
					clog.Warn(err)
				}
			}()
			c.Next()
		} else {
			clog.Warnf("limit = %d, service unavailable. url = %s", n, c.Request.RequestURI)
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}
}
