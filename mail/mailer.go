package mail

import (
	"bytes"
	"crypto/tls"
	"embed"
	"encoding/json"
	"github.com/twibber/core/cfg"
	"log/slog"
	"strconv"
	"text/template"

	"gopkg.in/gomail.v2"
)

// Embed HTML and text templates into the binary for easier deployment and management of these files.
//
//go:embed templates/html/*
var htmlTemplates embed.FS

//go:embed templates/text/*
var textTemplates embed.FS

var (
	mailer   *gomail.Dialer     // Mailer configuration for sending emails.
	textTmpl *template.Template // Compiled text templates for emails.
	htmlTmpl *template.Template // Compiled HTML templates for emails.
)

func init() {
	// Load and parse email templates.
	loadTemplates()

	// If the application is in debug mode, do not initialize the mailer.
	if cfg.Config.Debug {
		slog.Debug("mailer initialized in debug mode")
		return
	}

	// Initialize the mailer with configuration from the environment or config files.
	port, err := strconv.Atoi(cfg.Config.MailPort)
	if err != nil {
		slog.With(
			"error", err,
			"port", cfg.Config.MailPort,
		).Error("failed to parse mail port")
		panic(err)
	}

	// Set up mailer with TLS configuration based on application security requirements.
	mailer = gomail.NewDialer(cfg.Config.MailHost, port, cfg.Config.MailUsername, cfg.Config.MailPassword)
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: !cfg.Config.MailSecure, ServerName: cfg.Config.MailHost}

	slog.With("host", cfg.Config.MailHost,
		"port", cfg.Config.MailPort,
		"secure", cfg.Config.MailSecure,
		"username", cfg.Config.MailUsername,
		"sender", cfg.Config.MailSender,
		"reply", cfg.Config.MailReply,
	).Info("successfully configured")
}

// loadTemplates compiles the email templates from the embedded file system.
func loadTemplates() {
	var err error
	htmlTmpl, err = template.ParseFS(htmlTemplates, "templates/html/*")
	if err != nil {
		panic(err)
	}

	textTmpl, err = template.ParseFS(textTemplates, "templates/text/*")
	if err != nil {
		panic(err)
	}

	slog.With("html", "templates/html/*",
		"text", "templates/text/*",
	).Info("successfully loaded templates")
}

// Send constructs and sends an email using the specified subject, template, and data.
func Send(subject, templateName string, data interface{}) error {
	var htmlEmail, textEmail bytes.Buffer

	// Execute HTML template.
	if err := htmlTmpl.ExecuteTemplate(&htmlEmail, templateName+".html", data); err != nil {
		slog.With("template", templateName+".html").Error("failed to execute HTML template")
		return err
	}

	// Execute text template.
	if err := textTmpl.ExecuteTemplate(&textEmail, templateName+".txt", data); err != nil {
		slog.With("template", templateName+".txt").Error("failed to execute text template")
		return err
	}

	// Send the email.
	return sendEmail(subject, data, textEmail.String(), htmlEmail.String())
}

// sendEmail configures the email message and sends it.
func sendEmail(subject string, data interface{}, textContent, htmlContent string) error {
	var defaultData Defaults

	// Marshal and unmarshal the data to ensure it is in the correct format.
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.With("mailer",
			"data", data,
			"error", err,
		).Error("failed to marshal data")
		return err
	}

	// Unmarshal the data into the default data struct.
	if err := json.Unmarshal(jsonData, &defaultData); err != nil {
		slog.With("mailer",
			"data", data,
			"error", err,
		).Error("failed to unmarshal data")
		return err
	}

	// Create the email message.
	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.Config.MailSender)
	msg.SetAddressHeader("To", defaultData.Email, defaultData.Name)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", textContent)
	msg.AddAlternative("text/html", htmlContent)

	// In debug mode, log the email details instead of sending.
	if cfg.Config.Debug {
		slog.With(
			"from", cfg.Config.MailSender,
			"to", defaultData.Email,
			"subject", subject,
			"data", data,
		).Debug("email not sent (debug mode)")
		return nil
	}

	// Send the email.
	return mailer.DialAndSend(msg)
}
