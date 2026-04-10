package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/afifn11/gopay-x/services/wallet-service/internal/usecase"
)

type WalletHandler struct {
	walletUC usecase.WalletUsecase
	validate *validator.Validate
}

func NewWalletHandler(walletUC usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{
		walletUC: walletUC,
		validate: validator.New(),
	}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	wallet, err := h.walletUC.CreateWallet(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrWalletAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "wallet created", "data": wallet})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	wallet, err := h.walletUC.GetWallet(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": wallet})
}

func (h *WalletHandler) TopUp(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")

	var req usecase.TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userIDStr.(string)

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.walletUC.TopUp(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case usecase.ErrDuplicateTransaction:
			c.JSON(http.StatusOK, gin.H{"message": "duplicate request, returning existing transaction", "data": tx})
		case usecase.ErrInsufficientBalance:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case usecase.ErrWalletSuspended:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case usecase.ErrLockFailed:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "top up successful", "data": tx})
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.walletUC.GetTransactionHistory(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Internal endpoint — dipanggil oleh service lain
func (h *WalletHandler) CreateWalletInternal(c *gin.Context) {
	var body struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(body.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	wallet, err := h.walletUC.CreateWallet(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case usecase.ErrWalletAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "wallet created", "data": wallet})
}