package v1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/middleware"
	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/service"
)

var (
	testArticleFoo = domain.Article{
		ID:      1,
		UserID:  123,
		Title:   "title foo",
		Content: "content foo",
	}
	testArticleBar = domain.Article{
		ID:      2,
		UserID:  123,
		Title:   "title bar",
		Content: "content bar",
	}
	testArticles = []domain.Article{
		testArticleFoo,
		testArticleBar,
	}
	testErr = errors.New("something happened")
)

func TestArticleHandler_HandleCreateArticle(t *testing.T) {
	type fields struct {
		setupService func() ArticleService
	}
	type args struct {
		buildReqBody func() string
	}
	type want struct {
		article  *domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "201 Created",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockCreate = func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
						return &testArticleFoo, nil
					}
					return mock
				},
			},
			args: args{
				buildReqBody: func() string {
					article := request.CreateArticleRequest{
						UserID:  123,
						Title:   "title",
						Content: "content",
					}

					body, err := json.Marshal(article)
					require.NoError(t, err)

					return string(body)
				},
			},
			want: want{
				respCode: http.StatusCreated,
				article:  &testArticleFoo,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockCreate = func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
						return nil, testErr
					}
					return mock
				},
			},
			args: args{
				buildReqBody: func() string {
					article := request.CreateArticleRequest{
						UserID:  123,
						Title:   "title",
						Content: "content",
					}

					body, err := json.Marshal(article)
					require.NoError(t, err)

					return string(body)
				},
			},
			want: want{
				article:  nil,
				respCode: http.StatusInternalServerError,
				err:      response.NewInternalServerError(testErr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.fields.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.POST("/", h.HandleCreateArticle)

			// Prepare request.
			body := tt.args.buildReqBody()
			req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.want.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.want.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.want.err.ErrorCode, result.ErrorCode)
			} else {
				var result domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.article.UserID, result.UserID)
				assert.Equal(t, tt.want.article.Title, result.Title)
				assert.Equal(t, tt.want.article.Content, result.Content)
			}
		})
	}
}

func TestArticleHandler_HandleGetArticle(t *testing.T) {
	type fields struct {
		setupService func() ArticleService
	}
	type args struct {
		articleID string
	}
	type want struct {
		article  *domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "200 OK",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockGetArticle = func(ctx context.Context, id uint) (*domain.Article, error) {
						return &testArticleFoo, nil
					}
					return mock
				},
			},
			args: args{
				articleID: "999",
			},
			want: want{
				article:  &testArticleFoo,
				respCode: http.StatusOK,
				err:      nil,
			},
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockGetArticle = func(ctx context.Context, id uint) (*domain.Article, error) {
						return nil, testErr
					}
					return mock
				},
			},
			args: args{
				articleID: "999",
			},
			want: want{
				article:  nil,
				respCode: http.StatusInternalServerError,
				err:      response.NewInternalServerError(testErr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.fields.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.GET("/:articleID", h.HandleGetArticle)

			// Prepare request.
			req, err := http.NewRequest(http.MethodGet, "/"+tt.args.articleID, nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.want.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.want.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.want.err.ErrorCode, result.ErrorCode)
			} else {
				var result domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.article.UserID, result.UserID)
				assert.Equal(t, tt.want.article.Title, result.Title)
				assert.Equal(t, tt.want.article.Content, result.Content)
			}
		})
	}
}

func TestArticleHandler_HandleListArticles(t *testing.T) {
	type fields struct {
		setupService func() ArticleService
	}
	type want struct {
		articles []domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		wantErr bool
	}{
		{
			name: "200 OK",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockListArticles = func(ctx context.Context, per, perPage uint) ([]domain.Article, error) {
						return testArticles, nil
					}
					return mock
				},
			},
			want: want{
				articles: testArticles,
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockListArticles = func(ctx context.Context, per, perPage uint) ([]domain.Article, error) {
						return nil, testErr
					}
					return mock
				},
			},
			want: want{
				articles: nil,
				respCode: http.StatusInternalServerError,
				err:      response.NewInternalServerError(testErr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.fields.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.GET("/", middleware.Paginate(), h.HandleListArticles)

			// Prepare request.
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.want.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.want.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.want.err.ErrorCode, result.ErrorCode)
			} else {
				var result []domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, len(tt.want.articles), len(result))

				for i, v := range result {
					wantArticle := tt.want.articles[i]

					assert.Equal(t, wantArticle.UserID, v.UserID)
					assert.Equal(t, wantArticle.Title, v.Title)
					assert.Equal(t, wantArticle.Content, v.Content)
				}
			}
		})
	}
}

func TestArticleHandler_HandleListArticles_NotUsingPaginationMiddleware(t *testing.T) {
	// Prepare handler.
	svc := service.NewArticleServiceMock()
	h := NewArticleHandler(svc)

	// Create router and attach handler.
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/", h.HandleListArticles) // Without loading pagination middleware.

	// Prepare request.
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	// Execute request.
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	// Check the response code.
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var result response.ErrResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	want := response.NewInternalServerError(errors.New("something went wrong"))
	assert.Equal(t, want.StatusCode, result.StatusCode)
	assert.Equal(t, want.ErrorMsg, result.ErrorMsg)
	assert.Equal(t, want.ErrorCode, result.ErrorCode)
}

func TestArticleHandler_HandleSearchArticles(t *testing.T) {
	type fields struct {
		setupService func() ArticleService
	}
	type want struct {
		articles []domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		fields  fields
		want    want
		wantErr bool
	}{
		{
			name: "200 OK",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockSearchArticles = func(ctx context.Context, title, content string) ([]domain.Article, error) {
						return testArticles, nil
					}
					return mock
				},
			},
			want: want{
				articles: testArticles,
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "500 Internal Server Error - When service returns an error",
			fields: fields{
				setupService: func() ArticleService {
					mock := service.NewArticleServiceMock()
					mock.MockSearchArticles = func(ctx context.Context, title, content string) ([]domain.Article, error) {
						return nil, testErr
					}
					return mock
				},
			},
			want: want{
				articles: nil,
				respCode: http.StatusInternalServerError,
				err:      response.NewInternalServerError(testErr),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare handler.
			svc := tt.fields.setupService()
			h := NewArticleHandler(svc)

			// Create router and attach handler.
			gin.SetMode(gin.TestMode)
			r := gin.Default()
			r.GET("/", h.HandleSearchArticles)

			// Prepare request.
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			// Execute request.
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			// Check the response code.
			assert.Equal(t, tt.want.respCode, resp.Code)

			if tt.wantErr {
				var result response.ErrResponse
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, tt.want.err.StatusCode, result.StatusCode)
				assert.Equal(t, tt.want.err.ErrorMsg, result.ErrorMsg)
				assert.Equal(t, tt.want.err.ErrorCode, result.ErrorCode)
			} else {
				var result []domain.Article
				err := json.Unmarshal(resp.Body.Bytes(), &result)

				assert.NoError(t, err)
				assert.Equal(t, len(tt.want.articles), len(result))

				for i, v := range result {
					wantArticle := tt.want.articles[i]

					assert.Equal(t, wantArticle.UserID, v.UserID)
					assert.Equal(t, wantArticle.Title, v.Title)
					assert.Equal(t, wantArticle.Content, v.Content)
				}
			}
		})
	}
}
