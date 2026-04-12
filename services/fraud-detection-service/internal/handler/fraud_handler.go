package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/fraud-detection-service/internal/usecase"
)

type FraudHandler struct {
	fraudUC usecase.FraudUsecase
}

func NewFraudHandler(fraudUC usecase.FraudUsecase) *FraudHandler {
	return &FraudHandler{fraudUC: fraudUC}
}

func (h *FraudHandler) GetChecks(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	checks, total, err := h.fraudUC.GetChecksByUserID(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": checks, "total": total})
}

func (h *FraudHandler) GetRiskProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	profile, err := h.fraudUC.GetRiskProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}