package categories

import (
	"net/http"

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
