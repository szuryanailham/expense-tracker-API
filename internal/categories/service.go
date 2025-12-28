package categories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
)

type CategoryService interface {
	ListCategoriesByUser(ctx context.Context,userID pgtype.UUID) ([]repo.Category, error)
	
}

type svc struct {
	repo repo.Querier
}


func NewService(repo repo.Querier) CategoryService {
	return &svc{repo: repo}
}


func (s *svc)ListCategoriesByUser(ctx context.Context, userID pgtype.UUID)([]repo.Category, error){
	// take userid when user login
	return s.repo.ListCategoriesByUser(ctx, userID)
}