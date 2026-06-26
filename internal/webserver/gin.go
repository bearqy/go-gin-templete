package webserver

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go-gin-templete/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const shutdownTimeoutSecond = 10

func NewEngine(router func(gin.IRouter), accessLog io.Writer) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	g.Use(RequestID())
	g.Use(gin.CustomRecoveryWithWriter(accessLog, func(c *gin.Context, recovered any) {
		slog.Error("panic recovered", "request_id", c.GetString("request_id"), slog.Any("error", recovered))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "internal_error",
				"message": "internal server error",
			},
		})
	}))
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: accessLog}))

	g.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	g.Any("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	//g.Any("/_/setlevel/:level", func(c *gin.Context) {
	//	level := c.Param("level")
	//	oldLevel := logger.SetLevel(level)
	//	if oldLevel == "" {
	//		c.String(400, "error log level")
	//		return
	//	}
	//	c.String(http.StatusOK, oldLevel)
	//})

	router(&g.RouterGroup)
	return g
}

func Run(ctx context.Context, webConfig config.WebConfig, accessLog io.Writer, router func(gin.IRouter)) error {
	if accessLog == nil {
		accessLog = io.Discard
	}
	g := NewEngine(router, accessLog)

	server := &http.Server{
		Addr:         webConfig.Address,
		Handler:      g,
		ReadTimeout:  time.Duration(webConfig.ReadTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(webConfig.WriteTimeoutSecond) * time.Second,
		IdleTimeout:  time.Duration(webConfig.IdleTimeoutSecond) * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("应用启动成功", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen error", slog.Any("error", err))
			errCh <- err
			return
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}

	timeoutContext, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSecond*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutContext); err != nil {
		slog.Error("Server Shutdown", slog.Any("error", err))
		return err
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	slog.Info("Server Shutdown success")
	return nil
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func newRequestID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buf[:])
}
