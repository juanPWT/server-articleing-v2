package model

import "time"

type Category struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`
}

type CategoryRequest struct {
	Name string `json:"name" validate:"required,max=50"`
}
