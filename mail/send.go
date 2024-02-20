package mail

import "github.com/twibber/core/cfg"

// Defaults struct holds common fields for all email types.
type Defaults struct {
	Email string // Recipient's email address
	Name  string // Recipient's name
}

// VerifyDTO is a data structure for verification emails.
type VerifyDTO struct {
	Defaults
	Code string
}

// Send dispatches a verification email using predefined template and subject.
func (data VerifyDTO) Send() error {
	// The template name "user_verify" should match a template file name (without extension)
	return Send("Verify your "+cfg.Config.Name+" Account", "user_verify", data)
}
