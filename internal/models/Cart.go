package models

type CartItem struct {
	ProductId string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	UserId string     `json:"userId"`
	Items  []CartItem `json:"items"`
}
