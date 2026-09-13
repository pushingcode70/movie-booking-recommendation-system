package dto

type CreateBookingRequest struct {
	ShowID        uint   `json:"show_id" binding:"required"`
	SeatIDs       []uint `json:"seat_ids" binding:"required,min=1"`
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type BookingResponse struct {
	ID          uint    `json:"id"`
	ShowID      uint    `json:"show_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`

	RazorpayOrderID string `json:"razorpay_order_id"`
	RazorpayKeyID   string `json:"razorpay_key_id"`
}
