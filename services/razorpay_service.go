package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"movie-booking/config"

	"github.com/razorpay/razorpay-go"
)

type RazorpayService struct {
	client *razorpay.Client
}

// NewRazorpayService creates a new Razorpay client.//constructor
func NewRazorpayService() *RazorpayService {

	keyID := config.AppConfig.RazorpayKeyID
	keySecret := config.AppConfig.RazorpayKeySecret

	client := razorpay.NewClient(keyID, keySecret)

	return &RazorpayService{
		client: client,
	}
}

// CreateOrder creates a new Razorpay Order.
// Razorpay expects the smallest currency unit (paise),
func (s *RazorpayService) CreateOrder(amount float64) (map[string]interface{}, error) {

	amountInPaise := int(amount * 100)

	data := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": "INR",
		"receipt":  fmt.Sprintf("booking_%d", amountInPaise),
	}

	order, err := s.client.Order.Create(data, nil)

	if err != nil {
		return nil, err
	}

	return order, nil
}

// "Did this payment confirmation really come from Razorpay, or did someone fake it?"
func (s *RazorpayService) VerifyPaymentSignature(orderID, paymentID, signature string) bool {
	//razorpay generates  the signature using "order_id|payment_id"
	body := orderID + "|" + paymentID

	// Create an HMAC-SHA256 hasher using the Razorpay secret key
	h := hmac.New(sha256.New, []byte(config.AppConfig.RazorpayKeySecret))

	// Hash the message (order_id|payment_id)
	h.Write([]byte(body))

	// Generate the expected signature in hexadecimal format
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// Securely compare the expected signature with the one received from Razorpay
	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(signature),
	)
}
