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

func (s *EmailService) SendTicketEmail(data dto.TicketEmailData) error {

	// load and parse html email template
	tmpl, err := template.ParseFiles("templates/ticket_email.html")

	if err != nil {
		return err
	}

	// buffer to hold generated html
	var body bytes.Buffer

	// render template with booking data
	err = tmpl.Execute(&body, data)

	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	// set sender email address
	m.SetHeader("From", s.email)

	// set recipient email address
	m.SetHeader("To", data.ToEmail)

	// set email subject
	m.SetHeader("Subject", "🎬 Your Movie Ticket"+data.MovieTitle)

	// set body content as html
	m.SetBody("text/html", body.String())

	// configure smtp dialer
	d := gomail.NewDialer(
		s.host,
		s.port,
		s.email,
		s.password,
	)
	// connect and send email
	return d.DialAndSend(m)
}

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
