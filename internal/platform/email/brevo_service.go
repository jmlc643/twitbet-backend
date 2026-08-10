package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"

	brevo "github.com/getbrevo/brevo-go/lib"
	"github.com/jmlc643/twitbet-backend/internal/identity/domain/port"
	"github.com/jmlc643/twitbet-backend/internal/platform/config"
)

type BrevoService struct {
	client      *brevo.APIClient
	sender      string
	templateDir string
}

func NewBrevoService(cfg *config.Config, templateDir string) (port.EmailService, error) {
	if cfg.BrevoAPIKey == "" {
		return nil, fmt.Errorf("brevo API key is missing")
	}

	brevoCfg := brevo.NewConfiguration()
	brevoCfg.AddDefaultHeader("api-key", cfg.BrevoAPIKey)

	return &BrevoService{
		client:      brevo.NewAPIClient(brevoCfg),
		sender:      cfg.SMTPSender,
		templateDir: templateDir,
	}, nil
}

func (s *BrevoService) renderTemplate(filename string, data interface{}) (string, error) {
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

func (s *BrevoService) send(ctx context.Context, toEmail, subject, body string) error {
	sender := &brevo.SendSmtpEmailSender{
		Email: s.sender,
	}

	email := brevo.SendSmtpEmail{
		Sender:      sender,
		To:          []brevo.SendSmtpEmailTo{{Email: toEmail}},
		Subject:     subject,
		HtmlContent: body,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to send email via brevo: %w", err)
	}

	return nil
}

func (s *BrevoService) SendVerificationEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("verification.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.send(ctx, toEmail, "Verifica tu cuenta - TwitBet", body)
}

func (s *BrevoService) SendPasswordResetEmail(ctx context.Context, toEmail, otpCode string) error {
	data := struct {
		OTPCode string
	}{
		OTPCode: otpCode,
	}

	body, err := s.renderTemplate("reset_password.html", data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.send(ctx, toEmail, "Recuperación de contraseña - TwitBet", body)
}
