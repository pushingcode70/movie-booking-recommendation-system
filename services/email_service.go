package services

import (
	"bytes"
	"html/template"
	"movie-booking/config"

	"movie-booking/dto"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	host     string
	port     int
	email    string
	password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:     config.AppConfig.SMTPHost,
		port:     config.AppConfig.SMTPPort,
		email:    config.AppConfig.SMTPEmail,
		password: config.AppConfig.SMTPPassword,
	}
}

// SendTicketEmail generates an HTML email from a template
// and sends the movie ticket to the customer.
func (s *EmailService) SendTicketEmail(data dto.TicketEmailData) error {

	// Load and parse the HTML email template.
	// The template contains placeholders (e.g. {{.MovieName}})
	// that will be replaced with actual booking data.
	tmpl, err := template.ParseFiles("templates/ticket_email.html")

	if err != nil {
		return err
	}

	// Create a buffer to hold the generated HTML after
	// the template is executed.
	var body bytes.Buffer

	// Fill the template with the booking data and write
	// the final HTML into the buffer.
	err = tmpl.Execute(&body, data)

	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	// Set the sender's email address.
	m.SetHeader("From", s.email)

	// Set the recipient's email address.
	m.SetHeader("To", data.ToEmail)

	// Set the email subject.
	m.SetHeader("Subject", "🎬 Your Movie Ticket"+data.MovieTitle)

	// Set the email body as HTML.
	// body.String() converts the buffer into a string.
	m.SetBody("text/html", body.String())

	// Configure the SMTP server using the email credentials.
	d := gomail.NewDialer(
		s.host,
		s.port,
		s.email,
		s.password,
	)
	// Connect to the SMTP server and send the email.
	return d.DialAndSend(m)
}

// SendVerificationOTP sends an OTP to the user's email
// SendVerificationOTP sends an OTP email using an HTML template.
func (s *EmailService) SendVerificationOTP(data dto.OTPEmailData) error {

	tmpl, err := template.ParseFiles("templates/otp_email.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer

	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", s.email)
	m.SetHeader("To", data.ToEmail)
	m.SetHeader("Subject", data.Title)

	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(
		s.host,
		s.port,
		s.email,
		s.password,
	)

	return d.DialAndSend(m)
}
