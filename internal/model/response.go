package model

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Hata mesajı"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message" example:"İşlem başarılı"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status" example:"healthy"`
	Service string `json:"service" example:"url-shortener"`
	Version string `json:"version" example:"1.0.0"`
}
