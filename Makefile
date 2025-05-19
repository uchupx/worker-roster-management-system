
.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generated
	go mod tidy
	go mod vendor

config:
	@echo "Generating config..."
	mkdir -p ./config
	cp .env.example ./config/.env

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen -generate types,server,spec -o ./generated/api.gen.go -package api api.yml


