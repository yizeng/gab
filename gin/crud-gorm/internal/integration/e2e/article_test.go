package e2e

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/yizeng/gab/gin/crud-gorm/internal/api"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/request"
	"github.com/yizeng/gab/gin/crud-gorm/internal/api/handler/v1/response"
	"github.com/yizeng/gab/gin/crud-gorm/internal/config"
	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/gin/crud-gorm/pkg/dockertester"
)

var (
	testArticle999 = domain.Article{
		ID:      999,
		UserID:  123,
		Title:   "seeded title 999",
		Content: "seeded content 999",
	}
	testArticle888 = domain.Article{
		ID:      888,
		UserID:  123,
		Title:   "seeded title 888",
		Content: "seeded content 888",
	}
)

type ArticleHandlerTestSuite struct {
	suite.Suite

	db       *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	server   *api.Server
}

func (s *ArticleHandlerTestSuite) SetupSuite() {
	// Initialize container.
	dt := dockertester.InitPostgres()
	s.pool = dt.Pool
	s.resource = dt.Resource

	// Open connection.
	db, err := dockertester.OpenPostgres(dt.Resource, dt.HostPort)
	require.NoError(s.T(), err)

	s.db = db
}

func (s *ArticleHandlerTestSuite) TearDownSuite() {
	err := s.pool.Purge(s.resource) // Destroy the container.
	require.NoError(s.T(), err)
}

func (s *ArticleHandlerTestSuite) SetupTest() {
	// Run migrations.
	err := dao.InitTables(s.db)
	require.NoError(s.T(), err)

	// Seed database.
	script, err := os.ReadFile("../scripts/seed_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)

	// Create API server.
	s.server = api.NewServer(&config.AppConfig{
		API: &config.APIConfig{},
		Gin: &config.GinConfig{
			Mode: gin.TestMode,
		},
		Postgres: &config.PostgresConfig{},
	}, s.db)
}

func (s *ArticleHandlerTestSuite) TearDownTest() {
	s.deleteAllArticles()
}

func (s *ArticleHandlerTestSuite) deleteAllArticles() {
	script, err := os.ReadFile("../scripts/delete_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, new(ArticleHandlerTestSuite))
}

func (s *ArticleHandlerTestSuite) TestArticleHandler_HandleCreateArticle() {
	type args struct {
		buildReqBody func() string
	}
	type want struct {
		article  domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "201 Created",
			args: args{
				buildReqBody: func() string {
					article := request.CreateArticleRequest{
						UserID:  123,
						Title:   "title 1",
						Content: "content 1",
					}

					body, err := json.Marshal(article)
					require.NoError(s.T(), err)

					return string(body)
				},
			},
			want: want{
				article: domain.Article{
					UserID:  123,
					Title:   "title 1",
					Content: "content 1",
				},
				respCode: http.StatusCreated,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "400 Bad Request - Missing user_id, Title too long, Content too long",
			args: args{
				buildReqBody: func() string {
					article := request.CreateArticleRequest{
						Title:   uniuri.NewLen(200),
						Content: uniuri.NewLen(10000),
					}

					body, err := json.Marshal(article)
					require.NoError(s.T(), err)

					return string(body)
				},
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewBadRequest(errors.New("content: the length must be between 1 and 5000; title: the length must be between 1 and 128; user_id: cannot be blank.")),
			},
			wantErr: true,
		},
		{
			name: "400 Bad Request - invalid JSON",
			args: args{
				buildReqBody: func() string {
					return "["
				},
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewBadRequest(errors.New("unexpected EOF")),
			},
			wantErr: true,
		},
		{
			name: "400 Bad Request - Already exists",
			args: args{
				buildReqBody: func() string {
					article := request.CreateArticleRequest{
						UserID:  123,
						Title:   "seeded title 999",
						Content: "seeded content 999",
					}

					body, err := json.Marshal(article)
					require.NoError(s.T(), err)

					return string(body)
				},
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewBadRequest(errors.New("article already exists")),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			body := tt.args.buildReqBody()
			req, err := http.NewRequest("POST", "/api/v1/articles", strings.NewReader(body))
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlerTestSuite) TestArticleHandler_HandleGetArticle() {
	type args struct {
		articleID string
	}
	type want struct {
		article  domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "200 OK",
			args: args{
				articleID: "999",
			},
			want: want{
				article:  testArticle999,
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "404 Not Found - articleID is not found",
			args: args{
				articleID: "1",
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusNotFound,
				err:      response.NewNotFound("article", "ID", "1"),
			},
			wantErr: true,
		},
		{
			name: "404 Not Found - articleID is negative",
			args: args{
				articleID: "-1",
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusNotFound,
				err:      response.NewNotFound("article", "ID", "-1"),
			},
			wantErr: true,
		},
		{
			name: "400 Bad Request - invalid articleID",
			args: args{
				articleID: "abc",
			},
			want: want{
				article:  domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewInvalidInput("articleID", "abc"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			req, err := http.NewRequest("GET", "/api/v1/articles/"+tt.args.articleID, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlerTestSuite) TestArticleHandler_HandleListArticles() {
	type args struct {
		query string
	}
	type want struct {
		articles []domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		setup   func()
		args    args
		want    want
		wantErr bool
	}{
		{
			name:  "200 OK",
			setup: func() {},
			args: args{
				query: "",
			},
			want: want{
				articles: []domain.Article{
					testArticle999,
					testArticle888,
				},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name:  "200 OK - With pagination",
			setup: func() {},
			args: args{
				query: "?page=2&per_page=1",
			},
			want: want{
				articles: []domain.Article{
					testArticle888,
				},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "200 OK - When there are no articles",
			setup: func() {
				s.deleteAllArticles()
			},
			args: args{
				query: "?page=1&per_page=2",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name:  "400 Bad Request - Invalid page query",
			setup: func() {},
			args: args{
				query: "?page=abc&per_page=2",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewInvalidInput("page", "abc"),
			},
			wantErr: true,
		},
		{
			name:  "400 Bad Request - Negative page query",
			setup: func() {},
			args: args{
				query: "?page=-123&per_page=2",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewInvalidInput("page", "-123"),
			},
			wantErr: true,
		},
		{
			name:  "400 Bad Request - Invalid per_page query",
			setup: func() {},
			args: args{
				query: "?page=1&per_page=abc",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewInvalidInput("per_page", "abc"),
			},
			wantErr: true,
		},
		{
			name:  "400 Bad Request - Negative per_page query",
			setup: func() {},
			args: args{
				query: "?page=1&per_page=-123",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusBadRequest,
				err:      response.NewInvalidInput("per_page", "-123"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Setup the tests.
			tt.setup()

			// Prepare Request.
			req, err := http.NewRequest("GET", "/api/v1/articles"+tt.args.query, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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

func (s *ArticleHandlerTestSuite) TestArticleHandler_HandleSearchArticles() {
	type args struct {
		query string
	}
	type want struct {
		articles []domain.Article
		respCode int
		err      *response.ErrResponse
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "200 OK - by title",
			args: args{
				query: "title=999",
			},
			want: want{
				articles: []domain.Article{
					testArticle999,
				},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "200 OK - by content",
			args: args{
				query: "content=999",
			},
			want: want{
				articles: []domain.Article{
					testArticle999,
				},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "200 OK - When there are no results",
			args: args{
				query: "title=no-title&content=no-content",
			},
			want: want{
				articles: []domain.Article{},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
		{
			name: "200 OK - No query parameters",
			args: args{
				query: "",
			},
			want: want{
				articles: []domain.Article{
					testArticle999,
					testArticle888,
				},
				respCode: http.StatusOK,
				err:      nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Prepare Request.
			req, err := http.NewRequest("GET", "/api/v1/articles/search?"+tt.args.query, nil)
			require.NoError(t, err)

			// Execute Request.
			resp := executeRequest(req, s.server)

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
