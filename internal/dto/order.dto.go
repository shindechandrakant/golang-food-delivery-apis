package dto

type OrderItemRequest struct {
	ProductId string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gte=1"`
}

type PlaceOrderRequest struct {
	Items      []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
	CouponCode string             `json:"couponCode"`
}

type OrderItemResponse struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderResponse struct {
	Id       string              `json:"id"`
	Items    []OrderItemResponse `json:"items"`
	Products []ProductResponse   `json:"products"`
}
