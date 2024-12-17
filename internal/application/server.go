package application

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprom "github.com/zsais/go-gin-prometheus"
	"log/slog"
	"nft_service/infrastructure/config"
	"nft_service/infrastructure/database"
	"nft_service/infrastructure/rabbit"
	"nft_service/infrastructure/utils"
	"nft_service/internal/contract"
	"nft_service/internal/controller"
	"nft_service/internal/persistence"
	"nft_service/internal/service"
	"nft_service/internal/worker"
	"strings"
)

func setupServer(ctx context.Context, db *database.DB, cfg *config.Config, mq *rabbit.RabbitMQ) (*gin.Engine, error) {

	l := slog.Default()

	tokenRepo := persistence.NewTokenRepo(db.Conn)
	transferRepo := persistence.NewTransferRepo(db.Conn)

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

	go contractService.StartCacheUpdater(ctx, cfg.CacheUpdateInterval)

	tokenQueue, err := mq.DeclareQueue("token_queue")
	if err != nil {
		return nil, errors.New("failed to declare token queue" + err.Error())
	}

	transferQueue, err := mq.DeclareQueue("transfer_queue")
	if err != nil {
		return nil, errors.New("failed to declare transfer queue" + err.Error())
	}

	workerService, err := worker.NewWorker(contractUrl, mq, tokenQueue, transferQueue, tokenRepo, transferRepo, contractABI)
	if err != nil {
		return nil, errors.New("failed to create worker service" + err.Error())
	}

	go func() {
		if err := workerService.TokenUpdater(); err != nil {
			l.Error("failed to update token queue" + err.Error())
		}
	}()

	tokenService := service.NewTokenService(tokenRepo, contractService, mq, tokenQueue)
	transferService := service.NewTransferService(transferRepo, contractService, mq, transferQueue)
	tokenHandler := controller.NewTokenHandler(tokenService)
	transferHandler := controller.NewTransferHandler(transferService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(controller.LoggerMiddleware())

	prom := ginprom.NewPrometheus("gin")
	prom.Use(r)

	r.GET("/api/ping", controller.Ping)
	r.GET("/api/docs/spec", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	r.GET("/api/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/docs/spec")))

	r.POST("/api/tokens/create", tokenHandler.Create)
	r.GET("/api/tokens/list", tokenHandler.List)
	r.GET("/api/tokens/total_supply", tokenHandler.Total)
	r.GET("/api/tokens/total_supply_exact", tokenHandler.ExactTotal)

	r.POST("/api/transfers/create", transferHandler.Create)
	r.GET("/api/transfers/list", transferHandler.List)

	return r, nil
}
