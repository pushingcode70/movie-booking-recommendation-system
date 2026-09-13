package dto

type VerifyPaymentRequest struct {
	//This is exactly what Razorpay Checkout sends back after a successful payment.
	RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
	RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
	RazorpaySignature string `json:"razorpay_signature" binding:"required"`
}
