package db

import (
	"context"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository/dao"
	"github.com/yizeng/gab/chi/crud-gorm/pkg/dockertester"
)

var (
	hostPort string
	pool     *dockertest.Pool
	resource *dockertest.Resource
	repo     *repository.ArticleRepository
)

type ArticleDBTestSuite struct {
	suite.Suite

	db *gorm.DB
}

func (s *ArticleDBTestSuite) SetupSuite() {
	// Initialize container.
	hostPort, pool, resource = dockertester.InitDockertestForPostgres()

	// Open connection.
	db, err := dockertester.OpenPostgres(resource, hostPort)
	require.NoError(s.T(), err)

	s.db = db
}

func (s *ArticleDBTestSuite) TearDownSuite() {
	err := pool.Purge(resource) // Destroy the container.
	require.NoError(s.T(), err)
}

func (s *ArticleDBTestSuite) SetupTest() {
	// Run migrations.
	err := dao.InitTables(s.db)
	require.NoError(s.T(), err)

	// Seed database.
	script, err := os.ReadFile("../scripts/seed_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)

	// Initialize repository.
	daoArticle := dao.NewArticleDAO(s.db)
	repo = repository.NewArticleRepository(daoArticle)
}

func (s *ArticleDBTestSuite) TearDownTest() {
	s.deleteAllArticles()
}

func (s *ArticleDBTestSuite) deleteAllArticles() {
	script, err := os.ReadFile("../scripts/delete_articles.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)
}

func TestArticleDB(t *testing.T) {
	suite.Run(t, new(ArticleDBTestSuite))
}

func (s *ArticleDBTestSuite) TestArticleDB_FindByID() {
	result, err := repo.FindByID(context.TODO(), 999)
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.EqualValues(s.T(), result.UserID, 123)
	assert.Equal(s.T(), result.Title, "seeded title 999")
	assert.Equal(s.T(), result.Content, "seeded content 999")
}

func (s *ArticleDBTestSuite) TestArticleDB_FindAll() {
	result, err := repo.FindAll(context.TODO())
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), len(result), 2)
	assert.EqualValues(s.T(), result[0].UserID, 123)
	assert.Equal(s.T(), result[0].Title, "seeded title 999")
	assert.Equal(s.T(), result[0].Content, "seeded content 999")
}

func (s *ArticleDBTestSuite) TestArticleDB_Create() {
	result, err := repo.Create(context.TODO(), &domain.Article{
		UserID:  123,
		Title:   "new title",
		Content: "new content",
	})
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), result.Title, "new title")
	assert.Equal(s.T(), result.Content, "new content")
}
