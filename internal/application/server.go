package application

import (
	"errors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"nft_service/infrastructure/config"
	"nft_service/infrastructure/database"
	"nft_service/infrastructure/utils"
	"nft_service/internal/contract"
	"nft_service/internal/controller"
	"nft_service/internal/persistence"
	"nft_service/internal/service"
	"strings"
)

func setupServer(db *database.DB, cfg *config.Config) (*gin.Engine, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(controller.LoggerMiddleware())

	// ping
	r.GET("/api/ping", controller.Ping)

	// swagger docs
	r.GET("/api/docs/spec", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	r.GET("/api/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/docs/spec")))

	tokenRepo := persistence.NewTokenRepo(db.Conn)

	contractUrl, err := utils.GenerateInfuraURL(strings.ToLower(cfg.NetworkName), cfg.InfuraApiKey)
	if err != nil {
		return nil, errors.New("failed to generate Infura URL" + err.Error())
	}

	contractABI, err := utils.LoadABIFromFile(cfg.ContractABIPath)
	if err != nil {
		return nil, errors.New("failed to load contract ABI" + err.Error())
	}

	contractService, err := contract.NewNFTContract(contractUrl, cfg, contractABI)
	if err != nil {
		return nil, errors.New("failed to create contract service" + err.Error())
	}

	tokenService := service.NewTokenService(tokenRepo, contractService)
	tokenHandler := controller.NewTokenHandler(tokenService)

	r.POST("/api/tokens/create", tokenHandler.CreateToken)

	return r, nil
}
