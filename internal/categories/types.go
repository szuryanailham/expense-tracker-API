package categories

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCategoryRequest struct {
	UserID pgtype.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Type   string    `json:"type"`
}
type UpdateCategoryRequest struct {
	ID     pgtype.UUID     `json:"id"`     
	UserID pgtype.UUID   `json:"user_id"`  
	Name   string        `json:"name"`     
	Type   string        `json:"type"`     
}


type CreateCategoryResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}