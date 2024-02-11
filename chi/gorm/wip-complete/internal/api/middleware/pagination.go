package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/gorm/wip-complete/internal/api/handler/v1/response"
)

const (
	defaultPage    = 1
	defaultPerPage = 10
	maxPerPage     = 100
)

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageNumber := r.URL.Query().Get(PageQueryKey)
		parsedPageNumber, err := parsePageNumber(pageNumber)
		if err != nil {
			render.Render(w, r, response.ErrInvalidInput(PageQueryKey, pageNumber))

			return
		}

		perPage := r.URL.Query().Get(PerPageQueryKey)
		parsedPerPage, err := parsePerPage(perPage)
		if err != nil {
			render.Render(w, r, response.ErrInvalidInput(PerPageQueryKey, perPage))

			return
		}

		ctx := context.WithValue(r.Context(), PageQueryKey, parsedPageNumber)
		ctx = context.WithValue(ctx, PerPageQueryKey, parsedPerPage)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
