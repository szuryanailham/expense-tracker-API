package users

import (
	"log"
	"net/http"

	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
	"github.com/szuryanailham/expense-tracker/internal/json"
	"github.com/szuryanailham/expense-tracker/internal/middleware"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
    var payload RegisterRequest

    if err := json.Read(r, &payload); err != nil {
        log.Println(err)
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    result, err := h.service.Register(
        r.Context(),
        repo.CreateUserParams{
            FirstName: payload.FirstName,
            LastName:  payload.LastName,
            Email:     payload.Email,
            Password:  payload.Password,
        },
    )
    if err != nil {
        log.Println(err)
        http.Error(w, "failed to create new account", 
		http.StatusInternalServerError)
        return
    }
    json.Write(w, http.StatusCreated, map[string]string{
        "user_id": result.UserID,
        "token":   result.Token,
    })
}

func (h*handler)Login(w http.ResponseWriter, r *http.Request){
	   var payload LoginRequest
	   if err := json.Read(r, &payload); err != nil {
        log.Println(err)
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
	 result, err := h.service.Login(
        r.Context(),
        LoginParams{
            Email:    payload.Email,
            Password: payload.Password,
        },
    )
    if err != nil {
        log.Println(err)
        http.Error(w, "email atau password salah", http.StatusUnauthorized)
        return
    }
	json.Write(w, http.StatusOK, map[string]string{
        "user_id": result.UserID,
        "token":   result.Token,
    })
}

func ( h*handler)Authentication(w http.ResponseWriter, r *http.Request){
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, map[string]string{
		"user_id": userID,
		"message": "this is protected endpoint",
	})
}

