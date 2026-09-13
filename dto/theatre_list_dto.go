package dto

// theatreListResponse represents a theatre shown in customer theatre list
type TheatreListResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
