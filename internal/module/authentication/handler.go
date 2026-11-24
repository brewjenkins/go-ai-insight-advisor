package authentication

import (
	"errors"
	"log/slog"
	"net/http"

	"brewjenkins/go-ai-insight-advisor/internal/domain/constant"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc    *Service
	logger *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": constant.BAD_REQUEST})
		return
	}
	err := h.svc.Login(c.Request.Context(), req)
	if errors.Is(err, ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidCredentials.Error()})
		return
	}
	if err != nil {
		h.logger.Error("login failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
