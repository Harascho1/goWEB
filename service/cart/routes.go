package cart

import (
	"net/http"

	"github.com/Harascho1/goWEB/service/auth"
	"github.com/Harascho1/goWEB/types"
	"github.com/Harascho1/goWEB/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())
	var cartPayload types.CartCheckoutPayload
	err := utils.ParseJSON(r, &cartPayload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.Validator.Struct(cartPayload)
	if err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, error)
		return
	}

	productIDs, err := getCartItemsIDs(cartPayload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cartPayload.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"totalPrice": totalPrice,
		"orderId":    orderID,
	})
}
