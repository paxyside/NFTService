.PHONEY: all, format, pack, run, clear, docs
PATH = ./cmd/nft_service/main.go
NAME = ./nft_service
DOCS_PATH = ./docs/swagger.json


all: clear format	pack run  docs

format:
	go fmt ./...

pack:
	go build -o ${NAME} ${PATH}

run:
	${NAME}

clear:
	rm -f ${NAME}

docs:
	rm -f ${DOCS_PATH} && go-swagger3 --module-path . --main-file-path ${PATH} --output ${DOCS_PATH} --schema-without-pkg

