package controller

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"math/big"
	"net/http"
	"nft_service/internal/domain"
	"nft_service/internal/service"
	"strconv"
)

type TokenHandler struct {
	tokenService *service.TokenService
}

func NewTokenHandler(tokenService *service.TokenService) *TokenHandler {
	return &TokenHandler{tokenService: tokenService}
}

type CreateTokenRequest struct {
	Owner    string `json:"owner" binding:"required"`
	MediaUrl string `json:"media_url" binding:"required"`
}

type SupplyResponse struct {
	TotalSupply string `json:"total_supply"`
}

type ErrorResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
}

// Create NFT Token
// @Summary Create a new NFT token
// @Description Creates a new NFT token and assigns it to the provided owner's address with the specified media URL.
// @Tag NFT Token
// @Param token body CreateTokenRequest true "Data required to create the NFT token"
// @Success 201 {object} domain.Token "Successfully created token"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 500 {object} ErrorResponse "Failed to create token"
// @Router /api/tokens/create [post]
func (h *TokenHandler) Create(c *gin.Context) {

	var (
		l       = slog.Default()
		request = new(domain.Token)
	)

	if err := c.BindJSON(request); err != nil {
		l.Error("invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "invalid request",
		})
		return
	}

	token, err := h.tokenService.CreateToken(request)
	if err != nil {
		l.Error("failed to generate token", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusCreated, token)
}

// List Tokens
// @Summary Retrieve a paginated list of NFT tokens
// @Description Returns a list of NFT tokens. If `limit` and `offset` parameters are provided, they will be used for pagination. By default, `limit` is set to 200, and `offset` is 0. The `limit` value must be between 1 and 500.
// @Tag NFT Token
// @Param offset query int false "Pagination offset, default 0"
// @Param limit query int false "Number of pagination elements, default 200, max 500"
// @Success 200 {array} domain.Token "Successful response containing the list of tokens"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/tokens/list [get]
func (h *TokenHandler) List(c *gin.Context) {
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

	tokens, err := h.tokenService.ListTokens(limit, offset)
	if err != nil {
		l.Error("failed to list tokens", slog.Any("error", err))

		c.JSON(http.StatusInternalServerError, &gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to list tokens",
		})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Total
// @Summary Retrieve the total supply of NFT tokens from cache
// @Description Returns the total number of NFT tokens minted on the blockchain from cache.
// @Tag NFT Token
// @Success 200 {object} SupplyResponse "Successful response with total supply"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/tokens/total_supply [get]
func (h *TokenHandler) Total(c *gin.Context) {
	var (
		l           = slog.Default()
		totalSupply *big.Int
		err         error
	)

	c.Header("Content-Type", "application/json")

	totalSupply, err = h.tokenService.TotalSupply()
	if err != nil {
		l.Error("failed to get total supply", slog.Any("error", err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to get total supply",
		})
		return
	}

	response := SupplyResponse{
		TotalSupply: totalSupply.String(),
	}

	c.JSON(http.StatusOK, response)
}

// ExactTotal
// @Summary Retrieve the exact total supply of NFT tokens
// @Description Returns exact the total number of NFT tokens minted on the blockchain.
// @Tag NFT Token
// @Success 200 {object} SupplyResponse "Successful response with total supply"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/tokens/total_supply_exact [get]
func (h *TokenHandler) ExactTotal(c *gin.Context) {
	var (
		l           = slog.Default()
		totalSupply *big.Int
		err         error
	)

	c.Header("Content-Type", "application/json")

	totalSupply, err = h.tokenService.ExactTotalSupply()
	if err != nil {
		l.Error("failed to get total supply", slog.Any("error", err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"request_id": c.GetString("requestId"),
			"error":      "failed to get total supply",
		})
		return
	}

	response := SupplyResponse{
		TotalSupply: totalSupply.String(),
	}

	c.JSON(http.StatusOK, response)
}
