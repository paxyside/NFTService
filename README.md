# NFTService for Rock`n`Block
![img.png](img.png)

NFT Service is a backend application designed for managing ERC-721 tokens on the blockchain. \
The service enables minting unique NFTs using the mint method, retrieving the total supply of tokens, and listing stored NFTs and transfers from a PostgreSQL database.
It facilitates transferring tokens to new owners by creating and tracking blockchain transactions.\
The application integrates RabbitMQ queues to process token IDs once minted and to update transfer statuses after confirmation on the blockchain,
ensuring reliability and accurate record-keeping.

[Test task](https://confluence.rocknblock.io/pages/viewpage.action?pageId=1082566)

[Swagger UI](http://127.0.0.1:8008/api/docs/swagger/index.html)

[Prometheus UI](http://127.0.0.1:9090)

[RabbitMQ UI](http://127.0.0.1:15672)

## Installed Packages
- [Go-Ethereum](https://github.com/ethereum/go-ethereum)
- [GIN Web Framework](https://github.com/gin-gonic/gin)
- [Gin-Prometheus](https://github.com/zsais/go-gin-prometheus)
- [PGXPool for PostgreSQL](https://github.com/jackc/pgx)
- [Go-Swagger3](https://github.com/parvez3019/go-swagger3)
- [Go-Migrate](https://github.com/golang-migrate/migrate)
- [Go-RabbitMQ](https://github.com/rabbitmq/amqp091-go)

## Prerequisites
Step 1: Clone repository
```bash
git clone https://github.com/paxyside/NFTService.git nft_service
cd nft_service
```

Step 2: Environment Configuration
```bash
cp .env.example .env && chmod 600 .env
```

Step 3: Export Environment Variables
```bash
export $(cat .env | xargs)
```

Step 4: Create Docker Network (First Time Only)
```bash
docker network create ntf_network
```


## Start application using Docker Compose
```bash
docker compose build && docker compose up -d
```

## Start application using Makefile
Step 1: Start database container
```bash
docker compose up database -d
dokcer compose up rabbitmq -d
```
Step 2: Start application
```bash
make pack && make run
```

## Useful Commands

### To view logs use
```bash
docker compose logs --tail 100
```

### To update swagger docs use
```bash
rm -f ./docs/swagger.json &&
go-swagger3 --module-path . --main-file-path ./cmd/nft_service/main.go --output ./docs/swagger.json --schema-without-pkg
```


## Project structure
```bash
.
├── cmd/
│   └── nft_service/                 # Main entry point for the application
├── docs/
│   └── swagger.json                 # Swagger API documentation
├── http/
│   └── requests.http                # HTTP request examples for testing
├── infrastructure/
│   ├── config/                      # Application configuration (e.g., env parsing)
│   ├── database/                    # Database connection and initialization
│   ├── rabbit/                      # RabbitMQ connection and helpers
│   └── utils/                       # Utility functions (e.g., hashing, ABI loader)
├── internal/
│   ├── application/                 # Core application logic (e.g., server setup)
│   ├── contract/                    # Logic for interacting with blockchain contracts
│   ├── controller/                  # HTTP handlers, routing, and middleware
│   ├── domain/                      # Domain models and interfaces
│   ├── persistence/                 # Repositories and database interaction logic
│   ├── service/                     # Business services (e.g., token operations)
│   └── worker/                      # Asynchronous workers for blockchain updates
├── migrations/                      # Database migrations for schema
├── contract_abi.json                # ABI for smart contract interaction
├── docker-compose.yaml              # Docker Compose configuration
├── Dockerfile                       # Dockerfile for building the application
├── Makefile                         # Common build and run tasks
├── prometheus.yml                   # Prometheus monitoring configuration
├── go.mod                           # Go module dependencies
├── go.sum                           # Dependency checksums
├── README.md                        # Project documentation
└── img.png                          # Example or project illustration
```
