package controller

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"nft_service/internal/domain"
	"nft_service/internal/service"
	"strconv"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{transferService: transferService}
}

// Create
// @Summary Create transfer of NFT Token to new owner
// @Description Creates a new transfer of the NFT token to a new owner
// @Tag Transfers
// @Param token body CreateTransferRequest true "Data required to create the transfer NFT token"
// @Success 201 {object} domain.Transfer "Successfully created transfer"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 500 {object} ErrorResponse "Failed to create transfer"
// @Router /api/transfers/create [post]
func (h *TransferHandler) Create(c *gin.Context) {

	var (
		l       = slog.Default()
		request = new(domain.Transfer)
	)

	if err := c.BindJSON(request); err != nil {
		l.Error("invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid request",
		})
		return
	}

	token, err := h.transferService.CreateTransfer(request)
	if err != nil {
		l.Error("failed to generate transfer", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to generate transfer",
		})
		return
	}

	c.JSON(http.StatusCreated, token)
}

// List
// @Summary Retrieve a paginated list of transfers
// @Description Returns a list of transfers. If `limit` and `offset` parameters are provided, they will be used for pagination. By default, `limit` is set to 200, and `offset` is 0. The `limit` value must be between 1 and 500.
// @Tag Transfers
// @Param offset query int false "Pagination offset, default 0"
// @Param limit query int false "Number of pagination elements, default 200, max 500"
// @Success 200 {array} domain.Transfer "Successful response containing the list of transfers"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/transfers/list [get]
func (h *TransferHandler) List(c *gin.Context) {
	var l = slog.Default()

	c.Header("Content-Type", "application/json")

	limitQuery := c.DefaultQuery("limit", "200")
	offsetQuery := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		l.Error("atoi id", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid limit",
		})
		return
	}

	if limit < 1 || limit > 500 {
		l.Error("invalid limit", slog.Any("limit", limit), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid limit, must be between 1 and 500",
		})
		return
	}

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		l.Error("atoi id", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid offset",
		})
		return
	}

	if offset < 0 {
		l.Error("invalid offset", slog.Any("offset", offset), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid offset, must be greater than 0",
		})
		return
	}

	transfers, err := h.transferService.ListTransfer(limit, offset)
	if err != nil {
		l.Error("failed to list transfers", slog.Any("error", err))

		c.JSON(http.StatusInternalServerError, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to list transfers",
		})
		return
	}

	c.JSON(http.StatusOK, transfers)
}
