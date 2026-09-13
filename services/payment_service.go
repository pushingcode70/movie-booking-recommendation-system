package services

import (
	"errors"
	"log"
	"movie-booking/dto"
	"movie-booking/repositories"
	"strconv"

	"gorm.io/gorm"
)

type PaymentService struct {
	db              *gorm.DB
	repo            *repositories.PaymentRepository
	bookingRepo     *repositories.BookingRepository
	userRepo        *repositories.UserRepository
	showRepo        *repositories.ShowRepository
	movieRepo       *repositories.MovieRepository
	screenRepo      *repositories.ScreenRepository
	theatreRepo     *repositories.TheatreRepository
	bookingSeatRepo *repositories.BookingSeatRepository
	seatRepo        *repositories.SeatRepository
	razorpayService *RazorpayService
	emailService    *EmailService
}

// Constructor
func NewPaymentService(db *gorm.DB, repo *repositories.PaymentRepository,
	bookingRepo *repositories.BookingRepository,
	userRepo *repositories.UserRepository,
	showRepo *repositories.ShowRepository,
	movieRepo *repositories.MovieRepository,
	screenRepo *repositories.ScreenRepository,
	theatreRepo *repositories.TheatreRepository,
	bookingSeatRepo *repositories.BookingSeatRepository,
	seatRepo *repositories.SeatRepository,
	razorpayService *RazorpayService,
	emailService *EmailService) *PaymentService {
	return &PaymentService{
		db:              db,
		repo:            repo,
		bookingRepo:     bookingRepo,
		userRepo:        userRepo,
		showRepo:        showRepo,
		movieRepo:       movieRepo,
		screenRepo:      screenRepo,
		theatreRepo:     theatreRepo,
		bookingSeatRepo: bookingSeatRepo,
		seatRepo:        seatRepo,
		razorpayService: razorpayService,
		emailService:    emailService,
	}
}

// verify payment
func (s *PaymentService) VerifyPayment(req *dto.VerifyPaymentRequest) error {

	// 1. Find the payment record using the Razorpay Order ID
	// We saved this Order ID when creating the Razorpay order.
	payment, err := s.repo.GetPaymentByRazorpayOrderID(req.RazorpayOrderID)
	if err != nil {
		return err
	}

	if payment.Status == "SUCCESS" {
		return nil
	}

	// 2. Verify that the payment details really came from Razorpay
	// and were not tampered with by the client.
	valid := s.razorpayService.VerifyPaymentSignature(req.RazorpayOrderID, req.RazorpayPaymentID, req.RazorpaySignature)

	if !valid {
		return errors.New("invalid payment signature")

	}

	//store raozrpaypayment is id and mark payment as successful
	payment.RazorpayPaymentID = req.RazorpayPaymentID
	payment.Status = "SUCCESS"

	//update payment and booking together
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	paymentRepo := s.repo.WithTx(tx)
	bookingRepo := s.bookingRepo.WithTx(tx)

	if err := paymentRepo.UpdatePayment(payment); err != nil {
		tx.Rollback()
		return err
	}

	txBooking, err := bookingRepo.GetBookingByID(payment.BookingID)

	if err != nil {
		tx.Rollback()
		return err
	}

	if txBooking.Status == "CANCELLED" {
		tx.Rollback()
		return errors.New("booking is cancelled")
	}

	if err := bookingRepo.UpdateStatus(payment.BookingID, "CONFIRMED"); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Get booking
	booking, err := s.bookingRepo.GetBookingByID(payment.BookingID)
	if err != nil {
		return err
	}

	// Get user
	user, err := s.userRepo.GetUserByID(booking.UserID)
	if err != nil {
		return err
	}

	// Get show
	show, err := s.showRepo.GetShowByID(booking.ShowID)
	if err != nil {
		return err
	}

	// Get movie
	movie, err := s.movieRepo.GetMovieByID(show.MovieID)
	if err != nil {
		return err
	}

	// Get screen
	screen, err := s.screenRepo.GetScreenByID(show.ScreenID)
	if err != nil {
		return err
	}

	// Get theatre
	theatre, err := s.theatreRepo.GetTheatreByID(screen.TheatreID)
	if err != nil {
		return err
	}

	// Get all booked seats for this booking
	bookingSeats, err := s.bookingSeatRepo.GetBookingSeatsByBookingID(booking.ID)
	if err != nil {
		return err
	}

	// Extract only the Seat IDs from the booking-seat records.
	var seatIDs []uint

	for _, bs := range bookingSeats {
		seatIDs = append(seatIDs, bs.SeatID)
	}

	// Fetch complete seat details using the extracted Seat IDs.
	seats, err := s.seatRepo.GetSeatByIDs(seatIDs)
	if err != nil {
		return err
	}

	// Extract seat numbers (e.g., A1, A2, B3) for the ticket/email.
	var seatNumbers []string

	for _, seat := range seats {
		seatNumbers = append(seatNumbers, seat.SeatNumber)
	}

	ticketData := dto.TicketEmailData{
		ToEmail:      user.Email,
		CustomerName: user.Name,
		BookingID:    booking.ID,
		MovieTitle:   movie.Title,
		TheatreName:  theatre.Name,
		ScreenName:   strconv.Itoa(screen.ScreenNumber),
		ShowTime:     show.StartTime.Format("03:04 PM"),
		ShowDate:     show.StartTime.Format("02 Jan 2006"), //it converts into string as email dto takes string
		Seats:        seatNumbers,
		Amount:       booking.TotalAmount,
	}

	err = s.emailService.SendTicketEmail(ticketData)
	if err != nil {
		log.Printf("failed to send ticket email for booking %d: %v", booking.ID, err)
	}

	return nil
}
