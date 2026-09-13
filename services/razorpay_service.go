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

// constructor
func NewRazorpayService() *RazorpayService {

	keyID := config.AppConfig.RazorpayKeyID
	keySecret := config.AppConfig.RazorpayKeySecret

	client := razorpay.NewClient(keyID, keySecret)

	return &RazorpayService{
		client: client,
	}
}

// createOrder creates a new razorpay order
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

// verifyPaymentSignature checks if payment confirmation really came from razorpay
func (s *RazorpayService) VerifyPaymentSignature(orderID, paymentID, signature string) bool {
	// razorpay generates the signature using "order_id|payment_id"
	body := orderID + "|" + paymentID

	// create an hmac-sha256 hasher using the razorpay secret key
	h := hmac.New(sha256.New, []byte(config.AppConfig.RazorpayKeySecret))

	// hash the message (order_id|payment_id)
	h.Write([]byte(body))

	// generate expected signature in hexadecimal format
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	// securely compare expected signature with the one received from razorpay
	return hmac.Equal(
		[]byte(expectedSignature),
		[]byte(signature),
	)
}
