package categories

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
	"github.com/szuryanailham/expense-tracker/internal/json"
	"github.com/szuryanailham/expense-tracker/internal/middleware"
)

type handler struct {
	service CategoryService
}

func NewHandler(service CategoryService) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler)GetCategoriesById(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	categories, err := h.service.ListCategoriesByUser(r.Context(),userID)

	if err != nil {
		http.Error(w, "failed to fetch categories", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, categories)
}

func(h*handler)CreateCategory(w http.ResponseWriter, r *http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var tempCategory CreateCategoryRequest
	if err := json.Read(r, &tempCategory);err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	category, err := h.service.CreateCategory(r.Context(),repo.CreateCategoryParams{
		UserID: userID,
		Name: tempCategory.Name,
		Type: tempCategory.Type,
	})

	if err != nil {
		log.Println(err)
		http.Error(w,"Failed to create product", http.StatusInternalServerError)
		return
	}
	 json.Write(w, http.StatusCreated, CreateCategoryResponse{
        Status:  "success",
        Message: "Category created successfully",
        Data:    category,
    })
}

func(h* handler)UpdateCategory(w http.ResponseWriter, r *http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	categoryIDStr := chi.URLParam(r, "id")
	parsedUUID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}
	categoryID := pgtype.UUID{
    Bytes: parsedUUID,
    Valid: true,
}
	var tempCategory UpdateCategoryRequest
	if err := json.Read(r, &tempCategory);err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := h.service.UpdateCategory(r.Context(), repo.UpdateCategoryParams{
		ID:     categoryID,
		UserID: userID,          
		Name:   tempCategory.Name,
		Type:   tempCategory.Type,
	})

	if err != nil {
		http.Error(w, "failed to update category", http.StatusInternalServerError)
		return
	}

	 json.Write(w, http.StatusCreated, CreateCategoryResponse{
        Status:  "success",
        Message: "Category updated successfully",
        Data:    category,
    })

}
