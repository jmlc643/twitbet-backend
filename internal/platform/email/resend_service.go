package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/platform/config"
	"github.com/resend/resend-go/v2"
)

type ResendService struct {
	client      *resend.Client
	sender      string
	templateDir string
}

func NewResendService(cfg *config.Config, templateDir string) (port.EmailService, error) {
	if cfg.ResendAPIKey == "" {
		return nil, fmt.Errorf("resend API key is missing")
	}

	client := resend.NewClient(cfg.ResendAPIKey)

	return &ResendService{
		client:      client,
		sender:      cfg.SMTPSender,
		templateDir: templateDir,
	}, nil
}

func (s *ResendService) renderTemplate(filename string, data interface{}) (string, error) {
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

func (s *ResendService) SendVerificationEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("verification.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    s.sender,
		To:      []string{toEmail},
		Subject: "Verifica tu cuenta - TwitBet",
		Html:    body,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %w", err)
	}

	return nil
}

func (s *ResendService) SendPasswordResetEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("reset_password.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    s.sender,
		To:      []string{toEmail},
		Subject: "Recuperación de contraseña - TwitBet",
		Html:    body,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %w", err)
	}

	return nil
}
