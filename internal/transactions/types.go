package transactions

import "github.com/jackc/pgx/v5/pgtype"

type CreateTransactionRequest struct {
	UserID         	pgtype.UUID `json:"user_id"`
	CategoryID      pgtype.UUID `json:"category_id"`
	AmountCents     int64       `json:"amount_cents"`
	Note            pgtype.Text `json:"note"`
	TransactionDate pgtype.Date `json:"transaction_date"`
}

type TransactionsResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type UpdateTransactionRequest struct {
	UserID         	pgtype.UUID `json:"user_id"`
	CategoryID      pgtype.UUID `json:"category_id"`
	AmountCents     int64       `json:"amount_cents"`
	Note            pgtype.Text `json:"note"`
	TransactionDate pgtype.Date `json:"transaction_date"`
}