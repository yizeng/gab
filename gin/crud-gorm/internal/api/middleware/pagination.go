package middleware

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/response"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
	maxPerPage     = 100
)

func Paginate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageNumber := ctx.Query(PageQueryKey)
		parsedPageNumber, err := parsePageNumber(pageNumber)
		if err != nil {
			response.RenderErr(ctx, response.ErrInvalidInput(PageQueryKey, pageNumber))

			return
		}

		perPage := ctx.Query(PerPageQueryKey)
		parsedPerPage, err := parsePerPage(perPage)
		if err != nil {
			response.RenderErr(ctx, response.ErrInvalidInput(PerPageQueryKey, perPage))

			return
		}

		ctx.Set(PageQueryKey, parsedPageNumber)
		ctx.Set(PerPageQueryKey, parsedPerPage)

		ctx.Next()
	}
}

func parsePageNumber(number string) (uint, error) {
	if number == "" {
		return defaultPage, nil
	}

	parsed, err := strconv.Atoi(number)
	if err != nil || parsed <= 0 {
		return 0, errors.New("parse page query failed")
	}

	return uint(parsed), nil
}

func parsePerPage(perPage string) (uint, error) {
	if perPage == "" {
		return defaultPerPage, nil
	}

	parsed, err := strconv.Atoi(perPage)
	if err != nil || parsed <= 0 {
		return 0, errors.New("parse per page query failed")
	}

	if parsed > maxPerPage {
		return defaultPerPage, nil
	}

	return uint(parsed), nil
}
