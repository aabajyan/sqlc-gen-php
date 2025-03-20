.PHONY: build test bin/sqlc-gen-php bin/sqlc-gen-php.wasm

build:
	go build ./...

test:
	go test ./...

all: bin/sqlc-gen-php bin/sqlc-gen-php.wasm

bin/sqlc-gen-php: bin go.mod go.sum $(wildcard **/*.go)
	cd plugin && go build -o ../bin/sqlc-gen-php ./main.go

bin/sqlc-gen-php.wasm: bin/sqlc-gen-php
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-gen-php.wasm main.go

bin:
	mkdir -p bin

