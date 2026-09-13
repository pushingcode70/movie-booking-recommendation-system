package dto

type TicketEmailData struct {
	ToEmail      string
	CustomerName string

	MovieTitle  string
	TheatreName string
	ScreenName  string

	ShowDate string
	ShowTime string

	Seats []string

	BookingID uint
	Amount    float64
}
