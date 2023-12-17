all : run test coverage docker
.PHONY : all

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
