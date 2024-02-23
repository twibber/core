package models

import (
	"strings"
	"time"
)

type User struct {
	BaseModel

	Username    string `gorm:"size:64;not null;unique" json:"username"`
	DisplayName string `gorm:"size:512" json:"display_name"`

	Email string `gorm:"size:255;unique;not null" json:"email,omitempty"` // Ommitted for security reasons

	Connections []Connection `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"connections,omitempty"`

	Posts []Post `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// ConnectionType represents the type of connection.
type ConnectionType string

const (
	ProviderEmailType ConnectionType = "email"
	// For future use.
	// ProviderGoogleType ConnectionType = "google"
	// ProviderGitHubType ConnectionType = "github"
)

func (c ConnectionType) WithID(id string) string {
	return string(c) + ":" + id
}

// Connection represents an authentication connection, such as an email or an oauth connection.
type Connection struct {
	BaseModel

	// Email connection fields
	Password   string `gorm:"size:512" json:"-"`             // Password for the connection, only used for emails.
	Verified   bool   `gorm:"default:false" json:"verified"` // Whether the connection is verified, only used for emails.
	TOTPVerify string `gorm:"size:512" json:"-"`             // Time-based One-Time Password for verification of the email address

	// User owner of the connection
	UserID string `gorm:"not null" json:"-"`
	User   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// Sessions related to the connection
	Sessions []Session `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"sessions,omitempty"`
}

func (c *Connection) Type() ConnectionType {
	return ConnectionType(strings.Split(c.ID, ":")[0])
}

func (c *Connection) TypeID() ConnectionType {
	return ConnectionType(strings.Split(c.ID, ":")[1])
}

// Session represents an authenticated session related to a connection.
type Session struct {
	BaseModel

	ConnectionID string      `gorm:"not null" json:"-"`
	Connection   *Connection `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"connection,omitempty"`

	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}
