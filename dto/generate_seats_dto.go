package dto

type GenerateSeatsRequest struct {
	Rows        int `json:"rows" binding:"required,min=1,max=26"`
	SeatsPerRow int `json:"seats_per_row" binding:"required,min=1,max=50"`
}
