package transactions

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
	"github.com/szuryanailham/expense-tracker/internal/json"
	"github.com/szuryanailham/expense-tracker/internal/middleware"
)

type handler struct {
	service TransactionService
}

func NewHandler(service TransactionService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetCategoriesById(w http.ResponseWriter,r * http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	transactions, err := h.service.ListTransactionsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to fetch Transaction", http.StatusInternalServerError)
		return
	}

	json.Write(w,http.StatusOK,transactions)
}

func (h *handler)CreateTransaction(w http.ResponseWriter,r*http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	var tempTransaction CreateTransactionRequest
	if err := json.Read(r, &tempTransaction); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	Transactions, err := h.service.CreateTransaction(r.Context(),repo.CreateTransactionParams{
		UserID: userID,
		CategoryID: tempTransaction.CategoryID,
		AmountCents: tempTransaction.AmountCents,
		Note: tempTransaction.Note,
		TransactionDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil{
		log.Println(err)
		http.Error(w,"Failed to create Transactions", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, TransactionsResponse{
        Status:  "success",
        Message: "Category created successfully",
        Data:Transactions,
    })
}

func (h *handler)UpdateTransaction(w http.ResponseWriter,r*http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	transactionIDStr := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	transactionID := pgtype.UUID{
    Bytes: parsedUUID,
    Valid: true,
	}

	var tempTransaction UpdateTransactionRequest

	if err := json.Read(r, &tempTransaction); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

		Transactions, err := h.service.UpdateTransaction(r.Context(),repo.UpdateTransactionParams{
		UserID: userID,
		ID: transactionID,
		CategoryID: tempTransaction.CategoryID,
		AmountCents: tempTransaction.AmountCents,
		Note: tempTransaction.Note,
		TransactionDate: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
	})

		if err != nil{
		log.Println(err)
		http.Error(w,"Failed to update Transactions",
		 http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, TransactionsResponse{
        Status:  "success",
        Message: "Transaction updated successfully",
        Data:Transactions,
    })
	

}

func (h *handler)DeleteTransaction(w http.ResponseWriter,r*http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	transactionIDStr := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	transactionID := pgtype.UUID{
    Bytes: parsedUUID,
    Valid: true,
	}
	transaction, err := h.service.DeleteTransaction(r.Context(),repo.DeleteTransactionParams{
		ID:transactionID ,
		UserID: userID,
	})
	if err != nil {
		http.Error(w, "failed to delete category", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, TransactionsResponse{
        Status:  "success",
        Message: "Category deleted successfully.",
        Data:    transaction,
    })

}

func(h *handler)FindTransctionByID(w http.ResponseWriter,r *http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok{
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	transactionIDStr := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(transactionIDStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	transactionID := pgtype.UUID{
    Bytes: parsedUUID,
    Valid: true,
	}

	transaction, err := h.service.FindTransactionByID(r.Context(),repo.FindTransactionByIDParams{
		ID:transactionID ,
		UserID: userID,
	})
	if err != nil {
		http.Error(w, "Failed to retrieve transaction details", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, TransactionsResponse{
        Status:  "success",
        Message: "retrieve transaction successfully.",
        Data:    transaction,
    })
}

func (h *handler) GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	summary, err := h.service.GetMonthlySummary(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to retrieve monthly summary", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, TransactionsResponse{
		Status:  "success",
		Message: "Monthly summary retrieved successfully.",
		Data:    summary,
	})
}

func (h *handler) GetTransactionSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	summary, err := h.service.GetTransactionSummary(r.Context(), userID)
	if err != nil {
			http.Error(w, "Failed to retrieve transaction summary", http.StatusInternalServerError)
			return
		}

	json.Write(w, http.StatusOK, TransactionsResponse{
		Status:  "success",
		Message: "Transaction summary retrieved successfully.",
		Data:    summary,
	})
}