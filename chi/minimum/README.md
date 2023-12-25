# Go API Boilerplates - chi/minimum

This is a minimum example of Go API with Chi router.

A full list of libraries used can be found in [Dependencies](#dependencies) section.

## Development

This project can be booted up either

- natively with Go installed locally (assuming dependencies like PostgresSQL, Redis, etc. also exist)
- or by Docker Compose (PostgreSQL, Redis, etc. are included out of the box)

Please refer to [Makefile](./Makefile) or [docker-compose.yml](./docker-compose.yml) for details.

By default, the server will run at <http://localhost:3333>, but everything is configurable via [.env](.env) file.

### Tools

Some common development tools are introduced for better local development experience.

- [cosmtrek/air](https://github.com/cosmtrek/air) - Live reload for Go apps
- [GoTestTools/gotestfmt](https://github.com/GoTestTools/gotestfmt) - go test output for humans
- [swaggo/swag][swaggo/swag] - Automatically generate RESTful API documentation with Swagger 2.0 for Go.

```
go install github.com/cosmtrek/air@latest
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

Alternatively, use `make install` to install the required Go tooling locally.

### Documentation

This service has integrated [swaggo/swag][swaggo/swag] to create OpenAPI documentation automatically.

Use command `make swagg` or `swag init` to run the generation.

Then navigate to <http://localhost:3333/swagger/index.html> to view the API documentation.

## Dependencies

### API
- [go-chi/chi][go-chi/chi] - lightweight, idiomatic and composable router for building Go HTTP services
  - [go-chi/cors][go-chi/cors] - CORS net/http middleware for Go
  - [go-chi/render][go-chi/render] - easily manage HTTP request / response payloads of Go HTTP services
- [joho/godotenv][joho/godotenv] - A Go port of Ruby's dotenv library (Loads environment variables from .env files)
- [sethvargo/go-envconfig][sethvargo/go-envconfig] - A Go library for parsing struct tags from environment variables.
- [uber-go/zap][uber-go/zap] - Blazing fast, structured, leveled logging in Go.
- [swaggo/http-swagger][swaggo/http-swagger] - Default net/http wrapper to automatically generate RESTful API documentation with Swagger 2.0.

### Testing
- [stretchr/testify][stretchr/testify] - A toolkit with common assertions and mocks that plays nicely with the standard library

[go-chi/chi]: https://github.com/go-chi/chi
[go-chi/cors]: https://github.com/go-chi/cors
[go-chi/render]: https://github.com/go-chi/render
[joho/godotenv]: https://github.com/joho/godotenv
[sethvargo/go-envconfig]: https://github.com/sethvargo/go-envconfig
[uber-go/zap]: https://github.com/uber-go/zap
[stretchr/testify]: https://github.com/stretchr/testify
[swaggo/swag]: https://github.com/swaggo/swag
[swaggo/http-swagger]: https://github.com/swaggo/http-swagger
