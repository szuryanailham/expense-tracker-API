package categories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
)

type CategoryService interface {
	ListCategoriesByUser(ctx context.Context,userID pgtype.UUID) ([]repo.Category, error)
	CreateCategory(ctx context.Context, arg repo.CreateCategoryParams ) ([]repo.Category, error)
	UpdateCategory(ctx context.Context, arg repo.UpdateCategoryParams) ([]repo.Category, error)
}

type svc struct {
	repo repo.Querier
}


func NewService(repo repo.Querier) CategoryService {
	return &svc{repo: repo}
}


func (s *svc)ListCategoriesByUser(ctx context.Context, userID pgtype.UUID)([]repo.Category, error){
	return s.repo.ListCategoriesByUser(ctx, userID)
}

func (s *svc) CreateCategory(ctx context.Context, arg repo.CreateCategoryParams) ([]repo.Category, error) {
	category, err := s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		UserID: arg.UserID,
		Name:   arg.Name,
		Type:   arg.Type,
	})
	if err != nil {
		return []repo.Category{}, err
	}
		categorySlice := []repo.Category{category}
		return categorySlice,nil
}


func(s *svc)UpdateCategory(ctx context.Context, arg repo.UpdateCategoryParams)([]repo.Category, error){
	category, err := s.repo.UpdateCategory(ctx, repo.UpdateCategoryParams{
		Name: arg.Name,
		Type: arg.Type,
		ID: arg.ID,
		UserID: arg.UserID,
	})
	if err != nil {
		return []repo.Category{}, err
	}
	updatedCategory :=  []repo.Category{category}
	return updatedCategory,nil
}