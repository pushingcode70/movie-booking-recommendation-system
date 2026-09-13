package services

import (
	"errors"
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
	"os"

	"gorm.io/gorm"
)

type BookingService struct {
	db              *gorm.DB
	bookingRepo     *repositories.BookingRepository
	bookingSeatRepo *repositories.BookingSeatRepository
	paymentRepo     *repositories.PaymentRepository
	showRepo        *repositories.ShowRepository
	seatRepo        *repositories.SeatRepository
	razorpayService *RazorpayService
}

// constructor
func NewBookingService( //bcz booking needs all these things
	db *gorm.DB,
	bookingRepo *repositories.BookingRepository,
	bookingSeatRepo *repositories.BookingSeatRepository,
	paymentRepo *repositories.PaymentRepository,
	showRepo *repositories.ShowRepository,
	seatRepo *repositories.SeatRepository,
	razorpayService *RazorpayService,
) *BookingService {

	return &BookingService{
		db:              db,
		bookingRepo:     bookingRepo,
		bookingSeatRepo: bookingSeatRepo,
		paymentRepo:     paymentRepo,
		showRepo:        showRepo,
		seatRepo:        seatRepo,
		razorpayService: razorpayService,
	}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	return s.bookingRepo.CreateBooking(booking)
}

func (s *BookingService) GetBookingByID(id uint) (*models.Booking, error) {
	return s.bookingRepo.GetBookingByID(id)
}

func (s *BookingService) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	return s.bookingRepo.GetBookingsByUserID(userID)
}

func (s *BookingService) ConfirmBooking(id uint) error {
	return s.bookingRepo.UpdateStatus(id, "CONFIRMED")
}

func (s *BookingService) CancelBooking(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		bookingRepo := s.bookingRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)

		if err := bookingRepo.UpdateStatus(id, "CANCELLED"); err != nil {
			return err
		}

		payment, err := paymentRepo.GetPaymentByBookingID(id)
		if err != nil {
			return err
		}

		payment.Status = "CANCELLED"

		return paymentRepo.UpdatePayment(payment)
	})
}

func (s *BookingService) BookService(userID uint, req *dto.CreateBookingRequest) (*dto.BookingResponse, error) {

	var response *dto.BookingResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {

		// use transaction-aware repositories
		bookingRepo := s.bookingRepo.WithTx(tx)
		bookingSeatRepo := s.bookingSeatRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)
		showRepo := s.showRepo.WithTx(tx)
		seatRepo := s.seatRepo.WithTx(tx)

		// validate show
		show, err := showRepo.GetShowByID(req.ShowID)
		if err != nil {
			return err
		}

		// validate seats
		seats, err := seatRepo.GetSeatByIDs(req.SeatIDs)
		if err != nil {
			return err
		}

		// make sure every requested seat exists
		if len(seats) != len(req.SeatIDs) {
			return errors.New("one or more seats are not found")
		}

		// make sure every seat belongs to this show's screen
		for _, seat := range seats {
			if seat.ScreenID != show.ScreenID {
				return errors.New("seat belongs to another screen")
			}
		}

		// check whether any requested seat is already booked
		bookedSeatIDs, err := bookingSeatRepo.GetBookedSeatIDs(show.ID, req.SeatIDs)
		if err != nil {
			return err
		}

		if len(bookedSeatIDs) > 0 {
			return errors.New("one or more seats are already booked")
		}

		// calculate total amount
		totalAmount := show.Price * float64(len(seats))

		// create booking
		booking := &models.Booking{
			UserID:      userID,
			ShowID:      show.ID,
			TotalAmount: totalAmount,
			Status:      "PENDING",
		}

		err = bookingRepo.CreateBooking(booking)
		if err != nil {
			return err
		}

		// create bookingseat records
		for _, seatID := range req.SeatIDs {
			bookingSeat := &models.BookingSeat{
				BookingID: booking.ID,
				SeatID:    seatID,
			}

			err := bookingSeatRepo.CreateBookingSeat(bookingSeat)
			if err != nil {
				return err
			}
		}

		// create payment
		payment := &models.Payment{
			BookingID:     booking.ID,
			Amount:        totalAmount,
			PaymentMethod: req.PaymentMethod,
			Status:        "PENDING",
		}

		err = paymentRepo.CreatePayment(payment)
		if err != nil {
			return err
		}

		order, err := s.razorpayService.CreateOrder(totalAmount)
		if err != nil {
			return err
		}

		// extract the Razorpay Order ID from the response map.
		// the razorpay sdk returns a map[string]interface{}... so order["id"] has type interface
		// use a type assertion (.(string)) to convert it to a string before storing it.

		orderID, ok := order["id"].(string)
		if !ok || orderID == "" {
			return errors.New("invalid razorpay order id")
		}

		payment.RazorpayOrderID = orderID

		err = paymentRepo.UpdatePayment(payment)
		if err != nil {
			return err
		}

		// return response
		response = &dto.BookingResponse{
			ID:              booking.ID,
			ShowID:          booking.ShowID,
			TotalAmount:     booking.TotalAmount,
			Status:          booking.Status,
			RazorpayOrderID: payment.RazorpayOrderID,
			RazorpayKeyID:   os.Getenv("RAZORPAY_KEY_ID"),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
