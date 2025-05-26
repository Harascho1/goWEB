package product

import (
	"fmt"
	"net/http"

	"github.com/Harascho1/goWEB/types"
	"github.com/Harascho1/goWEB/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Hander struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Hander {
	return &Hander{store: store}
}

func (h *Hander) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handlerGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/add-product", h.handlerCreateProduct).Methods(http.MethodPost)
}

func (h *Hander) handlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.AddProductPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.Validator.Struct(payload)
	if err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, error)
		return
	}

	_, err = h.store.GetProductByName(payload.Name)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product with this name %s already exist", payload.Name))
		return
	}

	err = h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)

}

func (h *Hander) handlerGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ps)
}
