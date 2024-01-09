package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/yizeng/gab/chi/crud-gorm/internal/api/handler/v1/response"
)

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageNumber := r.URL.Query().Get(PageQueryKey)
		parsedPageNumber, err := parsePaginationQuery(pageNumber)
		if err != nil {
			render.Render(w, r, response.NewInvalidInput(PageQueryKey, pageNumber))
		}

		perPage := r.URL.Query().Get(PerPageQueryKey)
		parsedPerPage, err := parsePaginationQuery(perPage)
		if err != nil {
			render.Render(w, r, response.NewInvalidInput(PerPageQueryKey, perPage))
		}

		ctx := context.WithValue(r.Context(), PageQueryKey, parsedPageNumber)
		ctx = context.WithValue(ctx, PerPageQueryKey, parsedPerPage)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parsePaginationQuery(val string) (uint, error) {
	if val == "" {
		return 0, nil
	}

	parsed, err := strconv.Atoi(val)
	if err != nil || parsed < 0 {
		return 0, errors.New("parse pagination query failed")
	}

	return uint(parsed), nil
}
