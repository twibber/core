package models

// Post represents a post made by a user.
type Post struct {
	BaseModel

	// Author of the post
	AuthorID string `gorm:"not null" json:"author_id"`
	Author   *User  `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE" json:"author,omitempty"`

	Content string `gorm:"size:512" json:"content"`

	Likes []Like `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
}

type Like struct {
	BaseModel

	// Author of the post
	LikedByID string `gorm:"not null" json:"liked_by_id"`
	LikedBy   *User  `gorm:"foreignKey:LikedByID;references:ID;constraint:OnDelete:CASCADE" json:"liked_by,omitempty"`

	PostID string `gorm:"not null" json:"post_id"`
	Post   *Post  `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
}
