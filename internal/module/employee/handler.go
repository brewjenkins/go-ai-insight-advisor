package employee

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// EmployeeService is the consumer-side contract required by Handler.
type EmployeeService interface {
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*Employee, error)
	GetByID(ctx context.Context, id int64) (*Employee, error)
}

type Handler struct {
	svc    EmployeeService
	logger *slog.Logger
}

func NewHandler(svc EmployeeService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emp, err := h.svc.CreateEmployee(c.Request.Context(), req)
	if errors.Is(err, ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		h.logger.Error("create employee failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, toEmployeeResponse(emp))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	emp, err := h.svc.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}
	if err != nil {
		h.logger.Error("get employee failed", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, toEmployeeResponse(emp))
}
