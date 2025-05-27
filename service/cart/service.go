package cart

import (
	"fmt"

	"github.com/Harascho1/goWEB/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product id: %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}
	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userId int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	err := checkIfCartIsInStock(items, productMap)
	if err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil

}

func checkIfCartIsInStock(items []types.CartItem, products map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not avaliable in the store, please refresh a cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not avaliable in the quantity requested", product.Name)
		}
	}

	return nil
}
func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	var total float64

	for _, item := range items {
		product := productMap[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total

}
