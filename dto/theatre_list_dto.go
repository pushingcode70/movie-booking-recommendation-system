package dto

// TheatreListResponse represents a theatre shown
// in the customer theatre list.
type TheatreListResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
