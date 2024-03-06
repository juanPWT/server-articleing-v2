package model

import "time"

type Article struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	User_id      int       `json:"user_id" gorm:"not null"`
	Category_id  int       `json:"category_id" gorm:"not null"`
	Title        string    `json:"title" gorm:"type:varchar(255);not null"`
	Introduction string    `json:"introduction" gorm:"type:text;not null"`
	Thumbnail    string    `json:"thumbnail" gorm:"type:text;default:'https://placehold.co/400x400/png'"`
	IsPost       bool      `json:"is_post" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`
	User         User      `json:"user" gorm:"foreignKey:User_id;references:ID"`
	Category     Category  `json:"category" gorm:"foreignKey:Category_id;references:ID"`

	Body []Body
}

type Body struct {
	ID         int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Article_id int     `json:"article_id" gorm:"not null"`
	Content    string  `json:"content" gorm:"type:text;not null"`
	Article    Article `json:"article" gorm:"foreignKey:Article_id;references:ID"`
}

type CreteProject struct {
	User_id      int    `json:"user_id" validate:"required"`
	Category_id  int    `json:"category_id" validate:"required"`
	Title        string `json:"title" validate:"required"`
	Introduction string `json:"introduction" validate:"required"`
	Thumbnail    string `json:"thumbnail"`
}

type CreateContent struct {
	Content string `json:"content" validate:"required"`
}

type ArticleDetail struct {
	Article Article `json:"article"`
	Body    []Body  `json:"body"`
}
