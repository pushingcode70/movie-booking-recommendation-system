package dto

type CreateScreenRequest struct {
	TheatreID uint   `json:"theatre_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type UpdateScreenRequest struct {
	Name string `json:"name"`
}

type ScreenResponse struct {
	ID        uint   `json:"id"`
	TheatreID uint   `json:"theatre_id"`
	Name      string `json:"name"`
}
