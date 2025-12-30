package transactions

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
)

type TransactionService interface {
	ListTransactionsByUser(ctx context.Context, userID pgtype.UUID) ([]repo.ListTransactionsByUserRow, error)
	CreateTransaction(ctx context.Context, arg repo.CreateTransactionParams) (repo.Transaction, error)
	UpdateTransaction(ctx context.Context, arg repo.UpdateTransactionParams) (repo.Transaction, error)
	DeleteTransaction(ctx context.Context, arg repo.DeleteTransactionParams) (int64, error)
	FindTransactionByID(ctx context.Context, arg repo.FindTransactionByIDParams) (repo.FindTransactionByIDRow, error)
	GetMonthlySummary(ctx context.Context, userID pgtype.UUID) ([]repo.GetMonthlySummaryRow, error)
	GetTransactionSummary(ctx context.Context, userID pgtype.UUID) (repo.GetTransactionSummaryRow, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier)TransactionService{
	return &svc{repo:repo}
}

func (s *svc)ListTransactionsByUser(ctx context.Context,userID pgtype.UUID)([]repo.ListTransactionsByUserRow,error){
	return s.repo.ListTransactionsByUser(ctx,userID)
}

func (s *svc)CreateTransaction(ctx context.Context, arg repo.CreateTransactionParams)(repo.Transaction, error){
	transaction , err := s.repo.CreateTransaction(ctx,repo.CreateTransactionParams{
		UserID: arg.UserID,
		CategoryID: arg.CategoryID,
		AmountCents: arg.AmountCents,
		Note: arg.Note,
		TransactionDate: arg.TransactionDate,
	})

	if err != nil{
		return repo.Transaction{},err
	}
	return transaction,nil
}

func(s *svc) UpdateTransaction(ctx context.Context, arg repo.UpdateTransactionParams)(repo.Transaction, error){
	transaction, err := s.repo.UpdateTransaction(ctx, repo.UpdateTransactionParams{
		CategoryID: arg.CategoryID,
		AmountCents: arg.AmountCents,
		Note: arg.Note,
		TransactionDate: arg.TransactionDate,
		ID: arg.ID,
		UserID: arg.UserID,
	})
	if err != nil {
		return repo.Transaction{},err
	}
	return transaction, nil
}

func(s*svc)DeleteTransaction(ctx context.Context, arg repo.DeleteTransactionParams)(int64, error){
	 rowsAffected, err := s.repo.DeleteTransaction(ctx, repo.DeleteTransactionParams{
        ID:     arg.ID,
        UserID: arg.UserID,
    })
	 if err != nil {
        return 0, err
    }
	 if rowsAffected == 0 {
        return 0, fmt.Errorf("transaction not found or you do not have permission")
    }
	 return rowsAffected, nil
}

func (s*svc)FindTransactionByID(ctx context.Context,arg repo.FindTransactionByIDParams)(repo.FindTransactionByIDRow,error){
	transaction, err := s.repo.FindTransactionByID(ctx, repo.FindTransactionByIDParams{
		ID: arg.ID,
		UserID: arg.UserID,
	})
	 if err != nil {
        return repo.FindTransactionByIDRow{}, err
    }
	return transaction,nil
}

func (s *svc) GetMonthlySummary(ctx context.Context,userID pgtype.UUID) ([]repo.GetMonthlySummaryRow, error) {
	summary, err := s.repo.GetMonthlySummary(ctx, userID)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s*svc)GetTransactionSummary(ctx context.Context, userID pgtype.UUID)(repo.GetTransactionSummaryRow,error){
		summary, err := s.repo.GetTransactionSummary(ctx,userID)
			if err != nil {
		return repo.GetTransactionSummaryRow{}, err
	}
	return summary, nil
}




