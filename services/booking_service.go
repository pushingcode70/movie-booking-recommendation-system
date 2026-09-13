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

// Constructor
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

// Create Booking
func (s *BookingService) CreateBooking(booking *models.Booking) error {
	return s.bookingRepo.CreateBooking(booking)
}

// Get Booking By ID
func (s *BookingService) GetBookingByID(id uint) (*models.Booking, error) {
	return s.bookingRepo.GetBookingByID(id)
}

// Get Bookings By User
func (s *BookingService) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	return s.bookingRepo.GetBookingsByUserID(userID)
}

// confirm booking
func (s *BookingService) ConfirmBooking(id uint) error {
	return s.bookingRepo.UpdateStatus(id, "CONFIRMED")
}

// cancel booking
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

		// Use transaction-aware repositories
		bookingRepo := s.bookingRepo.WithTx(tx)
		bookingSeatRepo := s.bookingSeatRepo.WithTx(tx)
		paymentRepo := s.paymentRepo.WithTx(tx)
		showRepo := s.showRepo.WithTx(tx)
		seatRepo := s.seatRepo.WithTx(tx)

		//Validate Show
		show, err := showRepo.GetShowByID(req.ShowID) //req is a pointer to a CreateBookingRequest object.and it will access the DTO
		if err != nil {
			return err
		}

		//  Validate Seats
		seats, err := seatRepo.GetSeatByIDs(req.SeatIDs)
		if err != nil {
			return err
		}

		//Make sure every requested seat exists
		if len(seats) != len(req.SeatIDs) { //supose you request or client 1,2,3,4 ie req.seatIDs and db only has 1,2,3 so one returns lenth 4 and other 3 so not equal so its error
			return errors.New("one or more seats are not found")
		}

		// Make sure every seat belongs to this show's screen
		for _, seat := range seats {
			if seat.ScreenID != show.ScreenID {
				return errors.New("seat belongs to another screen")
			}
		}

		// Check whether any requested seat is already booked
		bookedSeatIDs, err := bookingSeatRepo.GetBookedSeatIDs(show.ID, req.SeatIDs)
		if err != nil {
			return err
		}

		if len(bookedSeatIDs) > 0 {
			return errors.New("one or more seats are already booked")
		}

		// 3. Calculate Total Amount
		totalAmount := show.Price * float64(len(seats))
		// 4. Create Booking
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

		// 5. Create BookingSeat records
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

		// 6. Create Payment
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

		// Extract the Razorpay Order ID from the response map.
		// The Razorpay SDK returns a map[string]interface{}, so order["id"] has type interface{}.
		// Use a type assertion (.(string)) to convert it to a string before storing it.

		orderID, ok := order["id"].(string)
		if !ok || orderID == "" {
			return errors.New("invalid razorpay order id")
		}

		payment.RazorpayOrderID = orderID

		err = paymentRepo.UpdatePayment(payment)
		if err != nil {
			return err
		}

		// 7. Return Response
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
