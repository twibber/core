package models

// Post represents a post made by a user.
type Post struct {
	BaseModel

	// Author of the post
	AuthorID string `gorm:"not null" json:"author_id"`
	Author   *User  `gorm:"foreignKey:AuthorID;references:ID;constraint:OnDelete:CASCADE" json:"author,omitempty"`

	Content string `gorm:"size:512" json:"content"`

	// Relations
	Likes []Like `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`

	// -- Replies
	// Parent is only used when a post is a reply to another post.
	ParentID *string `gorm:"null" json:"parent_id,omitempty"` // Parent is only used when a post is a reply to another post, therefore it is nullable and optional
	Parent   *Post   `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	// delete all replies when a post is deleted
	Replies []Post `gorm:"foreignKey:ParentID;references:ID;constraint:OnDelete:CASCADE" json:"replies,omitempty"` // delete all replies when a post is deleted
}

type Like struct {
	BaseModel

	// Author of the post
	LikedByID string `gorm:"not null" json:"liked_by_id"`
	LikedBy   *User  `gorm:"foreignKey:LikedByID;references:ID;constraint:OnDelete:CASCADE" json:"liked_by,omitempty"`

	PostID string `gorm:"not null" json:"post_id"`
	Post   *Post  `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
}
