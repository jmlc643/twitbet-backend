package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strconv"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/platform/config"
	"gopkg.in/gomail.v2"
)

type GomailService struct {
	dialer     *gomail.Dialer
	sender     string
	templateDir string
}

func NewGomailService(cfg *config.Config, templateDir string) (port.EmailService, error) {
	portInt, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp port: %w", err)
	}

	dialer := gomail.NewDialer(cfg.SMTPHost, portInt, cfg.SMTPUser, cfg.SMTPPass)

	return &GomailService{
		dialer:      dialer,
		sender:      cfg.SMTPSender,
		templateDir: templateDir,
	}, nil
}

func (s *GomailService) renderTemplate(filename string, data interface{}) (string, error) {
	path := filepath.Join(s.templateDir, filename)
	t, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *GomailService) SendVerificationEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("verification.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Verifica tu cuenta - TwitBet")
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *GomailService) SendPasswordResetEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("reset_password.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.sender)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Recuperación de contraseña - TwitBet")
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
