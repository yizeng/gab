all : run test coverage docker
.PHONY : all

run:
	@air

test:
	@go test ./... -v -race

coverage:
	@go test ./... -v -race -coverpkg=./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html coverage.out -o coverage.html
	@open coverage.html

docker:
	@docker compose up --build --force-recreate -V
