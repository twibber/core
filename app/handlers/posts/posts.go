package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
	"gorm.io/gorm"
	"time"
)

type PostForm struct {
	Content string `json:"content" validate:"required,max=512"`
}

// CreatePost handles the creation of new posts.
func CreatePost(c *fiber.Ctx) error {
	// Get the request body and validate it.
	var body PostForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Get the current session for our author.
	user := c.Locals("session").(models.Session)

	// Create the post.
	post := models.Post{
		AuthorID: user.Connection.UserID,
		Content:  body.Content,
	}
	if err := db.DB.Create(&post).Error; err != nil {
		return err
	}

	// Return the created post.
	return c.JSON(post)
}

func CreateReply(c *fiber.Ctx) error {
	// Get the request body and validate it.
	var body PostForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Get the post by its ID.
	var post models.Post
	if err := db.DB.Where(models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&post).Error; err != nil {
		return err
	}

	// Get the current session for our author.
	user := c.Locals("session").(models.Session)

	// Create the reply.
	reply := models.Post{
		AuthorID: user.Connection.UserID,
		Content:  body.Content,
		// Set the parent ID to the ID of the post we are replying to.
		ParentID: &post.ID,
	}

	// Create the reply in the database and return any errors.
	if err := db.DB.Create(&reply).Error; err != nil {
		return err
	}

	// Return the created post.
	return c.JSON(reply)
}

// ExtendedPost represents a post with its like counts and whether the current user liked the post.
type ExtendedPost struct {
	models.Post

	Liked  bool       `json:"liked"`  // Whether the current user liked the post.
	Counts PostCounts `json:"counts"` // The counts of the post.
}

type PostCounts struct {
	Likes   int64 `json:"likes"`   // Total likes on the post.
	Replies int64 `json:"replies"` // Total replies on the post.
}

// extendPost extends a post with its like counts and whether the current user liked the post.
func extendPost(post models.Post, userID string) ExtendedPost {
	// define the extended post
	extendedPost := ExtendedPost{
		Post:  post,
		Liked: false,
		Counts: PostCounts{
			Likes:   int64(len(post.Likes)),
			Replies: int64(len(post.Replies)),
		},
	}

	// Check if the current user is logged in
	if userID != "" {
		// Check if the current user liked the post.
		for _, like := range post.Likes {
			if like.LikedByID == userID {
				extendedPost.Liked = true
				break
			}
		}
	}

	return extendedPost
}

// ListPosts handles the retrieval of all posts with their like counts and whether the current user liked the post.
func ListPosts(c *fiber.Ctx) error {
	// Get all posts.
	var posts []models.Post
	if err := db.DB.
		Preload("Likes").
		Preload("Replies").
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		// only get top level posts
		Where("parent_id IS NULL").
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return err
	}

	// Get the id of the current user
	userID := utils.GetUserID(c)

	// Loop through the posts and extend them.
	var extendedPosts []ExtendedPost
	for _, post := range posts {
		extendedPosts = append(extendedPosts, extendPost(post, userID))
	}

	// Return the extended version of the posts.
	return c.JSON(extendedPosts)
}

// GetPost handles the retrieval of a single post by its ID.
func GetPost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.
		Preload("Likes").
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		Where(models.Post{
			BaseModel: models.BaseModel{ID: c.Params("post")},
		}).
		First(&post).Error; err != nil {
		return err
	}

	// Get the id of the current user
	userID := utils.GetUserID(c)

	return c.JSON(extendPost(post, userID))
}

// ListPostReplies handles the retrieval of all replies to a single post by its ID.
func ListPostReplies(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.
		Preload("Replies.Likes").
		Preload("Replies.Author", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		Where(models.Post{
			BaseModel: models.BaseModel{ID: c.Params("post")},
		}).
		First(&post).Error; err != nil {
		return err
	}

	// Get the id of the current user
	userID := utils.GetUserID(c)

	// Loop through the replies and extend them.
	var extendedReplies []ExtendedPost
	for _, reply := range post.Replies {
		extendedReplies = append(extendedReplies, extendPost(reply, userID))
	}

	// Return the extended version of the replies.
	return c.JSON(extendedReplies)
}

// DeletePost handles the deletion of a single post by its ID as long as the author is the one making the request, and it was created within the last 5 minutes.
func DeletePost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.Where(models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&post).Error; err != nil {
		return err
	}

	// Get the current session for our author.
	user := c.Locals("session").(models.Session)

	// Check if the author of the post is the same as the author of the session.
	if post.AuthorID != user.Connection.UserID {
		return utils.NewError(fiber.StatusForbidden, "You are not the author of this post.", nil)
	}

	// Check if the post was created within the last 5 minutes.
	if time.Since(post.CreatedAt) > 5*time.Minute {
		return utils.NewError(fiber.StatusForbidden, "You can only delete posts created within the last 5 minutes.", nil)
	}

	// Delete the post.
	if err := db.DB.Delete(&post).Error; err != nil {
		return err
	}

	// Return the deleted post.
	return c.SendStatus(fiber.StatusOK)
}
