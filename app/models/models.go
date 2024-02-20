package models

import (
	"time"

	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
)

// Models is a slice of all the models in the application.
var Models = []interface{}{
	&User{},
	&Connection{},
	&Session{},
}

// BaseModel defines the basic structure for database models.
type BaseModel struct {
	ID         string `gorm:"primaryKey" json:"id"` // ID is the primary key.
	Timestamps        // Timestamps for creation and update.
}

// Timestamps holds creation and update times.
type Timestamps struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Time of creation.
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Time of update.
}

// BeforeCreate is triggered before creating a new record.
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = utils.UUIDv4() // Assign a UUID if ID is empty.
	}
	return nil
}

// BeforeCreate sets timestamps before creating a record.
func (t *Timestamps) BeforeCreate(tx *gorm.DB) error {
	now := time.Now() // Current time.
	t.CreatedAt = now // Set creation time.
	t.UpdatedAt = now // Set update time for new record.
	return nil
}
