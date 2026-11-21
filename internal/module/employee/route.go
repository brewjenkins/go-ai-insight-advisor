package employee

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the employee module dependencies and registers HTTP routes.
func RegisterRoutes(r gin.IRouter, db *sql.DB) {
	repo := NewPostgresRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	employees := r.Group("/employees")
	employees.POST("", h.Create)
	employees.GET("/:id", h.GetByID)
}
