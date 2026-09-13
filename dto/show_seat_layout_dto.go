package dto

// final response returned to the customer
type ShowSeatLayoutResponse struct {
	ShowID       uint           `json:"show_id"`
	MovieTitle   string         `json:"movie_title"`
	TheatreName  string         `json:"theatre_name"`
	ScreenNumber int            `json:"screen_number"`
	Rows         []SeatRowGroup `json:"rows"`
}

// represents one row (a, b, c...)
type SeatRowGroup struct {
	Row   string     `json:"row"`
	Seats []SeatInfo `json:"seats"`
}

type SeatInfo struct {
	SeatID     uint   `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`
	Status     string `json:"status"`
}

type ShowSeatRow struct {
	ShowID       uint   `json:"show_id"`
	MovieTitle   string `json:"movie_title"`
	TheatreName  string `json:"theatre_name"`
	ScreenNumber int    `json:"screen_number"`

	SeatID     uint   `json:"seat_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`

	IsBooked bool `json:"is_booked"`
}
