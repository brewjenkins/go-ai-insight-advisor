package authentication

import (
	"log/slog"

	"brewjenkins/go-ai-insight-advisor/internal/module/employee"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r gin.IRouter, logger *slog.Logger, db *sqlx.DB) {
	h := NewHandler(NewService(employee.NewPostgresRepository(db)), logger)
	r.POST("/login", h.Login)
}
