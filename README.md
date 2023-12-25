# Gab (Go API Boilerplates)

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/yizeng/gab/graph/badge.svg?token=MIC6dQC41V)](https://codecov.io/gh/yizeng/gab)

Gab (Go API Boilerplates) is a set of boilerplates to spin up Go APIs quickly.

## How to Use?

All boilerplates are grouped by their API frameworks,
like `chi` for services using [go-chi/chi][go-chi/chi]
and `gin` for services built on [gin-gonic/gin][gin-gonic/gin].

Then each sub-folder contains a boilerplate that can spin up an API service.

For example, `minimum` folder contains a minimalistic Go service with 1 example API endpoint.
`crud-gorm` contains a basic CRUD API with [GORM][go-gorm/gorm].

- /chi
    - [/minimum/](./chi/minimum/README.md)
    - [/crud-gorm/](./chi/crud-gorm/README.md)
    - /crud-gorm-cached/
    - /crud-gorm-sqlc/
    - ...
- /gin
    - [/minimum/](./gin/minimum/README.md)
    - /crud-gorm/
    - ...
- ...

### How to Run?

Each boilerplate can be booted up either

- natively with Go installed locally (assuming dependencies like PostgresSQL, Redis, etc. also exist)
- or by Docker Compose (PostgreSQL, Redis, etc. are included out of the box)

Please refer to their `README.md` for details.

[go-chi/chi]: https://github.com/go-chi/chi
[gin-gonic/gin]: https://github.com/gin-gonic/gin
[go-gorm/gorm]: https://github.com/go-gorm/gorm