package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"brewjenkins/go-ai-insight-advisor/internal/config"
	"brewjenkins/go-ai-insight-advisor/internal/domain"
	"brewjenkins/go-ai-insight-advisor/internal/module/employee"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the Gin engine, attaches core middleware, and registers module routes.
func NewRouter(cfg *config.Config, db *sql.DB, logger *slog.Logger) *gin.Engine {
	if logger == nil {
		logger = slog.Default()
	}
	if cfg != nil && cfg.App.Env == domain.MODE_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestid.New())
	r.Use(requestLogger(logger))
	r.Use(corsMiddleware(cfg))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	employee.RegisterRoutes(r, db)

	return r
}

func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	corsCfg := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders: []string{"Content-Length", "X-Request-ID"},
		MaxAge:        12 * time.Hour,
	}

	var origins []string
	if cfg != nil {
		origins = cfg.App.CorsAllowedOrigins
	}
	if len(origins) == 0 || (len(origins) == 1 && origins[0] == "*") {
		// Wildcard origin cannot be combined with AllowCredentials=true in browsers/cors.
		corsCfg.AllowAllOrigins = true
		corsCfg.AllowCredentials = false
	} else {
		corsCfg.AllowOrigins = origins
		corsCfg.AllowCredentials = true
	}

	return cors.New(corsCfg)
}

func requestLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		logger.Info("http request",
			"request_id", requestid.Get(c),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
		)
	}
}
