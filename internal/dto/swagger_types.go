package dto

// Generic response wrappers used only for Swagger schema generation.

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"error message"`
}

// Typed wrappers so Swagger can show concrete data shapes.

type AuthSuccessResponse struct {
	Success bool         `json:"success" example:"true"`
	Data    AuthResponse `json:"data"`
}

type ProductListResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    []ProductResponse `json:"data"`
}

type ProductSuccessResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    ProductResponse `json:"data"`
}

type CartSuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    CartDetails `json:"data"`
}

type CartDetails struct {
	Cart           CartResponse `json:"cart"`
	TotalCartValue float64      `json:"totalCartValue" example:"29.99"`
}

type OrderSuccessResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    OrderResponse `json:"data"`
}

type IdempotencyKeyResponse struct {
	Success        bool   `json:"success" example:"true"`
	IdempotencyKey string `json:"idempotencyKey" example:"550e8400-e29b-41d4-a716-446655440000"`
	ExpiresIn      string `json:"expiresIn" example:"10 minutes"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Data    string `json:"data" example:"operation successful"`
}
