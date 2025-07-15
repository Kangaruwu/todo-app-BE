package service

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"go-backend-todo/internal/config"
)

// EmailService handles email operations
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendHTMLEmail(ctx context.Context, to, subject, htmlBody string) error
	SendVerificationEmail(ctx context.Context, to, username, token string) error
	SendPasswordResetEmail(ctx context.Context, to, username, token string) error
}

type emailService struct {
	cfg *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

// SendEmail sends a plain text email
func (s *emailService) SendEmail(ctx context.Context, to, subject, body string) error {
	if s.cfg.Email.SMTPUsername == "" {
		// Skip sending email if SMTP not configured (development mode)
		fmt.Printf("Email would be sent to %s: %s\n", to, subject)
		return nil
	}

	// Setup authentication
	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUsername, s.cfg.Email.SMTPPassword, s.cfg.Email.SMTPHost)

	// Compose message
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// Send email
	addr := fmt.Sprintf("%s:%d", s.cfg.Email.SMTPHost, s.cfg.Email.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.Email.FromEmail, []string{to}, []byte(msg))
}

// SendHTMLEmail sends an HTML email
func (s *emailService) SendHTMLEmail(ctx context.Context, to, subject, htmlBody string) error {
	if s.cfg.Email.SMTPUsername == "" {
		// Skip sending email if SMTP not configured (development mode)
		fmt.Printf("HTML Email would be sent to %s: %s\n", to, subject)
		return nil
	}

	// Setup authentication
	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUsername, s.cfg.Email.SMTPPassword, s.cfg.Email.SMTPHost)

	// Compose HTML message
	headers := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		s.cfg.Email.FromName, s.cfg.Email.FromEmail, to, subject)
	msg := headers + htmlBody

	// Send email
	addr := fmt.Sprintf("%s:%d", s.cfg.Email.SMTPHost, s.cfg.Email.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.Email.FromEmail, []string{to}, []byte(msg))
}

// SendVerificationEmail sends email verification email
func (s *emailService) SendVerificationEmail(ctx context.Context, to, username, token string) error {
    log.Println("Verification token:", token)
    return nil
    // subject := "Please verify your email address"

	// htmlBody := s.getVerificationEmailTemplate(username, token)

	// return s.SendHTMLEmail(ctx, to, subject, htmlBody)
}

// SendPasswordResetEmail sends password reset email
func (s *emailService) SendPasswordResetEmail(ctx context.Context, to, username, token string) error {
    log.Println("Password reset token:", token)
    return nil
	// subject := "Reset your password"

	// htmlBody := s.getPasswordResetEmailTemplate(username, token)

	// return s.SendHTMLEmail(ctx, to, subject, htmlBody)
}

// GetVerificationEmailTemplate returns HTML template for email verification
func (s *emailService) getVerificationEmailTemplate(username, token string) string {
	verificationURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/verify-email?token=%s", token)

	template := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Todo App!</h1>
        </div>
        <div class="content">
            <h2>Hello %s,</h2>
            <p>Thank you for registering with Todo App. To complete your registration and verify your email address, please click the button below:</p>
            <a href="%s" class="button">Verify Email Address</a>
            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p><a href="%s">%s</a></p>
            <p>This verification link will expire in 24 hours.</p>
            <p>If you didn't create an account with Todo App, please ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Todo App. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(template, username, verificationURL, verificationURL, verificationURL)
}

// GetPasswordResetEmailTemplate returns HTML template for password reset
// TODO: for testing purposes , the function name has been changed to GetPasswordResetEmailTemplate (capitalized first letter 'G')
func (s *emailService) getPasswordResetEmailTemplate(username, token string) string {
	resetURL := fmt.Sprintf("http://localhost:8080/api/v1/auth/reset-password?token=%s", token)

	template := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #FF6B6B; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .button { display: inline-block; padding: 12px 24px; background-color: #FF6B6B; color: white; text-decoration: none; border-radius: 4px; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 12px; margin-top: 20px; }
        .warning { background-color: #FFF3CD; border: 1px solid #FFEAA7; padding: 15px; border-radius: 4px; margin: 15px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hello %s,</h2>
            <p>We received a request to reset your password for your Todo App account. If you made this request, click the button below:</p>
            <a href="%s" class="button">Reset Password</a>
            <p>If the button doesn't work, you can also copy and paste this link into your browser:</p>
            <p><a href="%s">%s</a></p>
            <div class="warning">
                <strong>Security Notice:</strong>
                <ul>
                    <li>This link will expire in 1 hour</li>
                    <li>This link can only be used once</li>
                    <li>If you didn't request this reset, please ignore this email</li>
                </ul>
            </div>
            <p>If you didn't request a password reset, your account is still secure and no changes have been made.</p>
        </div>
        <div class="footer">
            <p>&copy; 2025 Todo App. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(template, username, resetURL, resetURL, resetURL)
}
