.PHONY: build test generate

all: bin/sqlc-gen-php.wasm

build:
	go build ./...

test:
	go test ./...

update_tests:
	go test -v ./internal/tests -update

bin/sqlc-gen-php.wasm:
	GOOS=wasip1 GOARCH=wasm go build -o ./bin/sqlc-gen-php.wasm ./plugin/main.go

sqlc.yaml: bin/sqlc-gen-php.wasm
	rm -f sqlc.yaml
	cp examples/minimal/sqlc.yaml sqlc.yaml
	sha256sum bin/sqlc-gen-php.wasm | awk '{print $$1}' | xargs -I {} sed -i 's/sha256: .*/sha256: {}/' sqlc.yaml
	sed -i "s|url: .*|url: file://$(PWD)/bin/sqlc-gen-php.wasm|" sqlc.yaml
	sed -i "s|sqlc/|examples/minimal/sqlc/|g" sqlc.yaml
	sed -i "s|src/Sqlc|examples/minimal/src/Sqlc|g" sqlc.yaml

generate: sqlc.yaml

bin:
	mkdir -p bin

