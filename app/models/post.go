package models

// Post represents a post made by a user.
type Post struct {
	BaseModel

	// Author of the post
	AuthorID string `gorm:"not null" json:"author_id"`
	Author   *User  `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE" json:"author,omitempty"`

	Content string `gorm:"size:512" json:"content"`
}
