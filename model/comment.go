package model

import "time"

type Comment struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Article_id int       `json:"article_id" gorm:"not null"`
	User_id    int       `json:"user_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Replies    int       `json:"replies" gorm:"type:int;default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`

	User    User    `json:"user" gorm:"foreignKey:User_id;references:ID"`
	Article Article `json:"article" gorm:"foreignKey:Article_id;references:ID"`
}

type Reply struct {
	ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Comment_id int       `json:"comment_id" gorm:"not null"`
	User_id    int       `json:"user_id" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`

	Comment Comment `json:"comment" gorm:"foreignKey:Comment_id;references:ID"`
	User    User    `json:"user" gorm:"foreignKey:User_id;references:ID"`
}

type CommentRequest struct {
	Content    string `json:"content" validate:"required"`
	Article_id int    `json:"article_id" validate:"required"`
	User_id    int    `json:"user_id" validate:"required"`
}

type ReplyRequest struct {
	Content    string `json:"content" validate:"required"`
	Comment_id int    `json:"comment_id" validate:"required"`
	User_id    int    `json:"user_id" validate:"required"`
}
