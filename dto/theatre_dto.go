package dto

type CreateTheatreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city" binding:"required"`
}

type UpdateTheatreRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}

type TheatreResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
}
