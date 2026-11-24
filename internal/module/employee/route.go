package employee

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// RegisterRoutes wires the employee module dependencies and registers HTTP routes.
func RegisterRoutes(r gin.IRouter, db *sqlx.DB, logger *slog.Logger) {
	repo := NewPostgresRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc, logger)

	employees := r.Group("/employees")
	employees.POST("", h.Create)
	employees.GET("/:id", h.GetByID)
}
