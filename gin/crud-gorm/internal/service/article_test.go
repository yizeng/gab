package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/yizeng/gab/gin/crud-gorm/internal/domain"
	"github.com/yizeng/gab/gin/crud-gorm/internal/repository"
)

func TestArticleService_CreateArticle(t *testing.T) {
	testArticle := &domain.Article{
		UserID:  999,
		Title:   "title 999",
		Content: "content 999",
	}
	testErr := errors.New("something happened")

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
				article: testArticle,
			},
			want:       testArticle,
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
				article: testArticle,
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
	testArticle := &domain.Article{
		UserID:  999,
		Title:   "title 999",
		Content: "content 999",
	}
	testErr := errors.New("something happened")

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
						return testArticle, nil
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				id:  999,
			},
			want:       testArticle,
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
				id:  testArticle.ID,
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
	testArticles := []domain.Article{
		{
			UserID:  999,
			Title:   "title 999",
			Content: "content 999",
		},
		{
			UserID:  888,
			Title:   "title 888",
			Content: "content 888",
		},
	}
	testErr := errors.New("something happened")

	type fields struct {
		repo ArticleRepository
	}
	type args struct {
		ctx      context.Context
		articles []domain.Article
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
					MockFindAll: func(ctx context.Context) ([]domain.Article, error) {
						return testArticles, nil
					},
				},
			},
			args: args{
				ctx:      context.TODO(),
				articles: testArticles,
			},
			want:       testArticles,
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Error Path",
			fields: fields{
				repo: &repository.ArticleRepositoryMock{
					MockFindAll: func(ctx context.Context) ([]domain.Article, error) {
						return nil, testErr
					},
				},
			},
			args: args{
				ctx:      context.TODO(),
				articles: testArticles,
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
			got, err := s.ListArticles(tt.args.ctx)
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
