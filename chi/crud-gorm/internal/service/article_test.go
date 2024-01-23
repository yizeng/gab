package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yizeng/gab/chi/crud-gorm/internal/domain"
	"github.com/yizeng/gab/chi/crud-gorm/internal/repository"
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

func TestArticleService_CreateArticle(t *testing.T) {
	type fields struct {
		repo ArticleRepository
	}
	type args struct {
		ctx     context.Context
		article *domain.Article
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *domain.Article
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Happy Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockCreate: func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
						return article, nil
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				article: &testArticleFoo,
			},
			want:       &testArticleFoo,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Error Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockCreate: func(ctx context.Context, article *domain.Article) (*domain.Article, error) {
						return nil, testErr
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				article: &testArticleFoo,
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "s.repo.Create -> something happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleService{
				repo: tt.fields.repo,
			}
			got, err := s.CreateArticle(tt.args.ctx, tt.args.article)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("CreateArticle() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateArticle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_GetArticle(t *testing.T) {
	type fields struct {
		repo ArticleRepository
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       *domain.Article
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Happy Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockFindByID: func(ctx context.Context, id uint) (*domain.Article, error) {
						return &testArticleFoo, nil
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  999,
			},
			want:       &testArticleFoo,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Error Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockFindByID: func(ctx context.Context, id uint) (*domain.Article, error) {
						return nil, testErr
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  testArticleFoo.ID,
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "s.repo.FindByID -> something happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleService{
				repo: tt.fields.repo,
			}
			got, err := s.GetArticle(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("FindByID() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_ListArticles(t *testing.T) {
	type fields struct {
		repo ArticleRepository
	}
	type args struct {
		ctx     context.Context
		page    uint
		perPage uint
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []domain.Article
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Happy Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockFindAll: func(ctx context.Context, page, perPage uint) ([]domain.Article, error) {
						return testArticles, nil
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				page:    1,
				perPage: 1,
			},
			want:       testArticles,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Error Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockFindAll: func(ctx context.Context, page, perPage uint) ([]domain.Article, error) {
						return nil, testErr
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				page:    1,
				perPage: 1,
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "s.repo.FindAll -> something happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleService{
				repo: tt.fields.repo,
			}
			got, err := s.ListArticles(tt.args.ctx, 0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("ListArticles() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListArticles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArticleService_SearchArticles(t *testing.T) {
	type fields struct {
		repo ArticleRepository
	}
	type args struct {
		ctx     context.Context
		title   string
		content string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       []domain.Article
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Happy Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockSearch: func(ctx context.Context, title, content string) ([]domain.Article, error) {
						return testArticles, nil
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				title:   "test title",
				content: "test content",
			},
			want:       testArticles,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Error Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockSearch: func(ctx context.Context, title, content string) ([]domain.Article, error) {
						return nil, testErr
					},
				},
			},
			args: args{
				ctx:     context.TODO(),
				title:   "test title",
				content: "test content",
			},
			want:       nil,
			wantErr:    true,
			wantErrMsg: "s.repo.Search -> something happened",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ArticleService{
				repo: tt.fields.repo,
			}
			got, err := s.SearchArticles(tt.args.ctx, "999", "")
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr && err.Error() != tt.wantErrMsg {
				t.Errorf("SearchArticles() errorMsg = %v, wantErrMsg %v", err.Error(), tt.wantErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchArticles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
