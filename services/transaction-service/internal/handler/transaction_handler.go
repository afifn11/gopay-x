package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/transaction-service/internal/usecase"
)

type TransactionHandler struct {
	txUC usecase.TransactionUsecase
}

func NewTransactionHandler(txUC usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{txUC: txUC}
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	txID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	tx, err := h.txUC.GetTransaction(c.Request.Context(), txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tx})
}

func (h *TransactionHandler) GetHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	req := &usecase.HistoryRequest{
		Type:   c.Query("type"),
		Status: c.Query("status"),
		From:   c.Query("from"),
		To:     c.Query("to"),
	}

	if page := c.Query("page"); page != "" {
		fmt.Sscanf(page, "%d", &req.Page)
	}
	if limit := c.Query("limit"); limit != "" {
		fmt.Sscanf(limit, "%d", &req.Limit)
	}

	result, err := h.txUC.GetHistory(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *TransactionHandler) GetSummary(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	summary, err := h.txUC.GetSummary(c.Request.Context(), userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}