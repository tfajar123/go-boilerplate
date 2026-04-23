package mailerService

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"go-boilerplate/apps/internal/config"
)

type MailerService struct {
	cfg config.MailConfig
}

func NewMailerService(cfg config.MailConfig) *MailerService {
	return &MailerService{cfg: cfg}
}

func (s *MailerService) ResetPasswordURL() string {
	return s.cfg.ResetPasswordURL
}

func (s *MailerService) SendResetPasswordEmail(toEmail, toName, resetLink string, expiresAt time.Time) error {
	if err := s.validateConfig(); err != nil {
		return err
	}

	from := mail.Address{
		Name:    s.cfg.FromName,
		Address: s.cfg.FromEmail,
	}

	to := mail.Address{
		Name:    toName,
		Address: toEmail,
	}

	subject := "Reset Password"
	htmlBody := fmt.Sprintf(`
<html>
  <body style="font-family: Arial, sans-serif; line-height: 1.6; color: #222;">
    <h2>Reset Password</h2>
    <p>Halo %s,</p>
    <p>Kami menerima permintaan untuk mereset password akun Anda.</p>
    <p>
      <a href="%s" style="display: inline-block; padding: 12px 20px; background: #1f6feb; color: #fff; text-decoration: none; border-radius: 6px;">
        Reset Password
      </a>
    </p>
    <p>Link ini akan kedaluwarsa pada %s.</p>
    <p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
  </body>
</html>`, toName, resetLink, expiresAt.Format(time.RFC1123))

	message := strings.Join([]string{
		fmt.Sprintf("From: %s", from.String()),
		fmt.Sprintf("To: %s", to.String()),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		`Content-Type: text/html; charset="UTF-8"`,
		"",
		htmlBody,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)

	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{toEmail}, []byte(message))
}

func (s *MailerService) validateConfig() error {
	switch {
	case s.cfg.Host == "":
		return fmt.Errorf("MAIL_HOST is required")
	case s.cfg.Port == 0:
		return fmt.Errorf("MAIL_PORT is required")
	case s.cfg.Username == "":
		return fmt.Errorf("MAIL_USERNAME is required")
	case s.cfg.Password == "":
		return fmt.Errorf("MAIL_PASSWORD is required")
	case s.cfg.FromEmail == "":
		return fmt.Errorf("MAIL_FROM_EMAIL is required")
	case s.cfg.ResetPasswordURL == "":
		return fmt.Errorf("MAIL_RESET_PASSWORD_URL is required")
	default:
		return nil
	}
}
