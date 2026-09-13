package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) WithTx(tx *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: tx,
	}
}

func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) GetPaymentByRazorpayOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.
		Where("razorpay_order_id = ?", orderID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepository) GetPaymentByBookingID(bookingID uint) (*models.Payment, error) {
	var payment models.Payment

	err := r.db.
		Where("booking_id = ?", bookingID).
		First(&payment).Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}
