.PHONEY: all, format, pack, run, clear
PATH = ./cmd/nft_service/main.go
NAME = ./nft_service


all: clear format	pack run

format:
	go fmt ./...

pack:
	go build -o ${NAME} ${PATH}

run:
	${NAME}

clear:
	rm -f ${NAME}
