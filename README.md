# Gab (Go API Boilerplates)

Gab is a set of Go API Boilerplates (Gab) to spin up Go APIs quickly.

## How to Use?

Gab (Go API Boilerplates) are grouped by their API frameworks,
like `chi` for services using [go-chi/chi][go-chi/chi]
and `gin` for services built on [gin-gonic/gin][gin-gonic/gin].

Then each sub-folder contains a boilerplate that can spin up an API service.
For example, `minimum` folder contains a minimalistic Go service with 1 example API endpoint.

```
/chi
    /minimum/
    /crud-gorm/
    /crud-gorm-cached/
    /crud-gorm-sqlc/
    ...
/gin
    /minimum/
    /crud-gorm/
    ...
...
```

### Development Tools

Some common development tools are introduced for better development experience.

- [cosmtrek/air](https://github.com/cosmtrek/air) - Live reload for Go apps
- [GoTestTools/gotestfmt](https://github.com/GoTestTools/gotestfmt) - go test output for humans

```
go install github.com/cosmtrek/air@latest
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
```

## `chi/minimum`

### API
- [go-chi/chi][go-chi/chi] - lightweight, idiomatic and composable router for building Go HTTP services
  - [go-chi/render][go-chi/render] - easily manage HTTP request / response payloads of Go HTTP services
- [spf13/viper][spf13/viper] - Go configuration with fangs
- [uber-go/zap][uber-go/zap] - Blazing fast, structured, leveled logging in Go.

### Testing
- [stretchr/testify][stretchr/testify] - A toolkit with common assertions and mocks that plays nicely with the standard library

## `chi/crud-gorm`

### API
- [go-chi/chi][go-chi/chi] - lightweight, idiomatic and composable router for building Go HTTP services
  - [go-chi/render][go-chi/render] - easily manage HTTP request / response payloads of Go HTTP services
- [spf13/viper][spf13/viper] - Go configuration with fangs
- [uber-go/zap][uber-go/zap] - Blazing fast, structured, leveled logging in Go.
- [go-ozzo/ozzo-validation][go-ozzo/ozzo-validation] - An idiomatic Go (golang) validation package. Supports configurable and extensible validation rules (validators) using normal language constructs instead of error-prone struct tags.

### Database
- [PostgreSQL][PostgreSQL] - The World's Most Advanced Open Source Relational Database
- [go-gorm/gorm][go-gorm/gorm] - The fantastic ORM library for Golang, aims to be developer friendly

### Testing
- [stretchr/testify][stretchr/testify] - A toolkit with common assertions and mocks that plays nicely with the standard library
- [ory/dockertest][ory/dockertest] - Write better integration tests! Dockertest helps you boot up ephermal docker images for your Go tests with minimal work.

[go-chi/chi]: https://github.com/go-chi/chi
[gin-gonic/gin]: https://github.com/gin-gonic/gin
[go-chi/render]: https://github.com/go-chi/render
[spf13/viper]: https://github.com/spf13/viper
[uber-go/zap]: https://github.com/uber-go/zap
[stretchr/testify]: https://github.com/stretchr/testify
[go-ozzo/ozzo-validation]: https://github.com/go-ozzo/ozzo-validation
[PostgreSQL]: https://www.postgresql.org/
[go-gorm/gorm]: https://github.com/go-gorm/gorm
[ory/dockertest]: https://github.com/ory/dockertest
