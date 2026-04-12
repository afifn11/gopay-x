package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/afifn11/gopay-x/services/audit-service/internal/usecase"
)

type AuditHandler struct {
	auditUC usecase.AuditUsecase
}

func NewAuditHandler(auditUC usecase.AuditUsecase) *AuditHandler {
	return &AuditHandler{auditUC: auditUC}
}

func (h *AuditHandler) GetByActor(c *gin.Context) {
	actorID := c.Param("actor_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, total, err := h.auditUC.GetByActorID(c.Request.Context(), actorID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs, "total": total})
}

func (h *AuditHandler) GetByResource(c *gin.Context) {
	resourceID := c.Param("resource_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	logs, total, err := h.auditUC.GetByResourceID(c.Request.Context(), resourceID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs, "total": total})
}

func (h *AuditHandler) QueryLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	req := &usecase.QueryRequest{
		Service: c.Query("service"),
		Event:   c.Query("event"),
		From:    c.Query("from"),
		To:      c.Query("to"),
		Page:    page,
		Limit:   limit,
	}

	logs, total, err := h.auditUC.QueryLogs(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": logs, "total": total})
}