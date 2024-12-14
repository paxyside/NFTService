package controller

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"nft_service/internal/domain"
	"nft_service/internal/service"
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

type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateToken
// @Summary Create NFT token
// @Description Creates new token to owner address
// @Tag NFT Token
// @Param token body CreateTokenRequest true "Data to create NFT token"
// @Success 201 object domain.Token "Token created"
// @Failure 400 object ErrorResponse "Invalid request"
// @Failure 500 object ErrorResponse "Failed to create tokens"
// @Router /api/tokens/create [post]
func (h *TokenHandler) CreateToken(c *gin.Context) {

	var (
		l       = slog.Default()
		request = new(domain.Token)
	)

	if err := c.BindJSON(request); err != nil {
		l.Error("invalid request", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := h.tokenService.CreateToken(request)
	if err != nil {
		l.Error("failed to generate token", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, token)
}
