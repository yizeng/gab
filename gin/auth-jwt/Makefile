all : install run test coverage docker swag
.PHONY : all

install:
	@go version
	@go install github.com/cosmtrek/air@latest
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

run:
	@air

test:
	@set -euo pipefail
	@go test ./... -json -v -race 2>&1 | tee /tmp/gotest.log | gotestfmt

coverage:
	@set -euo pipefail
	@go test ./... -json -v -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic 2>&1 | tee /tmp/gotest.log | gotestfmt
	@go tool cover -html coverage.out -o coverage.html
	@open coverage.html

docker:
	@docker compose up --build --force-recreate -V

swag:
	@swag init
