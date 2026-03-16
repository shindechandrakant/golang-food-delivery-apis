package dto

type AddCartItemRequest struct {
	ProductId string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=1"`
}

type CartItemResponse struct {
	ProductID string  `json:"productID"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CartResponse struct {
	UserID string             `json:"userId"`
	Items  []CartItemResponse `json:"items"`
}
