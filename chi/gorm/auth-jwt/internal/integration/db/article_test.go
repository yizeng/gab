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

	"github.com/yizeng/gab/chi/gorm/auth-jwt/internal/repository/dao"
	"github.com/yizeng/gab/chi/gorm/auth-jwt/pkg/dockertester"
)

type ArticleDBTestSuite struct {
	suite.Suite

	db       *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource

	articleDAO *dao.ArticleDAO
}

func (s *ArticleDBTestSuite) SetupSuite() {
	// Initialize container.
	dt := dockertester.InitPostgres()
	s.pool = dt.Pool
	s.resource = dt.Resource

	// Open connection.
	db, err := dockertester.OpenPostgres(dt.Resource, dt.HostPort)
	require.NoError(s.T(), err)

	s.db = db
}

func (s *ArticleDBTestSuite) TearDownSuite() {
	err := s.pool.Purge(s.resource) // Destroy the container.
	require.NoError(s.T(), err)
}

func (s *ArticleDBTestSuite) SetupTest() {
	// Run migrations.
	err := dao.InitTables(s.db)
	require.NoError(s.T(), err)

	// Seed database.
	script, err := os.ReadFile("../scripts/seed_db.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)

	// Initialize repository.
	s.articleDAO = dao.NewArticleDAO(s.db)
}

func (s *ArticleDBTestSuite) TearDownTest() {
	s.cleanDB()
}

func (s *ArticleDBTestSuite) cleanDB() {
	script, err := os.ReadFile("../scripts/clean_db.sql")
	require.NoError(s.T(), err)

	err = s.db.Exec(string(script)).Error
	require.NoError(s.T(), err)
}

func TestArticleDB(t *testing.T) {
	suite.Run(t, new(ArticleDBTestSuite))
}

func (s *ArticleDBTestSuite) TestArticleDB_FindByID() {
	result, err := s.articleDAO.FindByID(context.TODO(), 99999)
	assert.Error(s.T(), gorm.ErrRecordNotFound)

	result, err = s.articleDAO.FindByID(context.TODO(), 999)
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.EqualValues(s.T(), 123, result.UserID)
	assert.Equal(s.T(), "seeded title 999", result.Title)
	assert.Equal(s.T(), "seeded content 999", result.Content)
}

func (s *ArticleDBTestSuite) TestArticleDB_FindAll() {
	result, err := s.articleDAO.FindAll(context.TODO(), 1, 10)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), len(result), 2)
	assert.EqualValues(s.T(), 123, result[0].UserID)
	assert.Equal(s.T(), "seeded title 999", result[0].Title)
	assert.Equal(s.T(), "seeded content 999", result[0].Content)

	result, err = s.articleDAO.FindAll(context.TODO(), 2, 1)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), 1, len(result))
	assert.EqualValues(s.T(), 123, result[0].UserID)
	assert.Equal(s.T(), "seeded title 888", result[0].Title)
	assert.Equal(s.T(), "seeded content 888", result[0].Content)
}

func (s *ArticleDBTestSuite) TestArticleDB_FindAll_Empty() {
	s.cleanDB()

	result, err := s.articleDAO.FindAll(context.TODO(), 1, 10)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
}

func (s *ArticleDBTestSuite) TestArticleDB_Insert() {
	result, err := s.articleDAO.Insert(context.TODO(), dao.Article{
		UserID:  999,
		Title:   "seeded title 99",
		Content: "duplicated",
	})
	assert.Error(s.T(), dao.ErrArticleDuplicated)

	result, err = s.articleDAO.Insert(context.TODO(), dao.Article{
		UserID:  123,
		Title:   "new title",
		Content: "new content",
	})
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "new title", result.Title)
	assert.Equal(s.T(), "new content", result.Content)
}

func (s *ArticleDBTestSuite) TestArticleDB_Search() {
	result, err := s.articleDAO.Search(context.TODO(), "999", "")
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.EqualValues(s.T(), 123, result[0].UserID)
	assert.Equal(s.T(), "seeded title 999", result[0].Title)
	assert.Equal(s.T(), "seeded content 999", result[0].Content)

	result, err = s.articleDAO.Search(context.TODO(), "", "999")
	assert.NoError(s.T(), err)

	assert.NotNil(s.T(), result)
	assert.EqualValues(s.T(), 123, result[0].UserID)
	assert.Equal(s.T(), "seeded title 999", result[0].Title)
	assert.Equal(s.T(), "seeded content 999", result[0].Content)
}
