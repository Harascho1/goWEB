package user

import (
	"fmt"
	"net/http"

	"github.com/Harascho1/goWEB/service/auth"
	"github.com/Harascho1/goWEB/types"
	"github.com/Harascho1/goWEB/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handlerLogin).Methods("POST")
	router.HandleFunc("/register", h.handlerRegister).Methods("POST")
}

func (h *Handler) handlerLogin(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handlerRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with this email '%s' already exist", payload.Email))
		return
	}

	hashPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
