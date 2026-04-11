package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/payment-service/internal/usecase"
)

type PaymentHandler struct {
	paymentUC usecase.PaymentUsecase
	validate  *validator.Validate
}

func NewPaymentHandler(paymentUC usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUC: paymentUC,
		validate:  validator.New(),
	}
}

func (h *PaymentHandler) Transfer(c *gin.Context) {
	senderIDStr, _ := c.Get("user_id")

	var req usecase.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.SenderUserID = senderIDStr.(string)

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentUC.CreateTransfer(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrDuplicatePayment:
			c.JSON(http.StatusOK, gin.H{"message": "duplicate request", "data": payment})
		case usecase.ErrSameUser:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case usecase.ErrLockFailed:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "transfer successful", "data": payment})
}

func (h *PaymentHandler) TopUpViaGateway(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")

	var req usecase.GatewayTopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userIDStr.(string)

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.paymentUC.CreateTopUpViaGateway(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrDuplicatePayment:
			c.JSON(http.StatusOK, gin.H{"message": "duplicate request", "data": payment})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "payment created", "data": payment})
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	payment, err := h.paymentUC.GetPaymentByID(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (h *PaymentHandler) GetHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.paymentUC.GetPaymentHistory(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *PaymentHandler) HandleCallback(c *gin.Context) {
	var req usecase.CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.paymentUC.HandleCallback(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "callback processed"})
}