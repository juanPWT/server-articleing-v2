package model

import "time"

type User struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username          string    `json:"username" gorm:"type:varchar(50);not null"`
	Email             string    `json:"email" gorm:"unique;not null"`
	Verification_code string    `json:"verification_code" gorm:"type:varchar(255);default:null"`
	Verified_email    bool      `json:"verified_email" gorm:"default:false"`
	Password          string    `json:"password" gorm:"not null"`
	Image             string    `json:"image" gorm:"type:text;default:'https://placehold.co/400x400/png'"`
	Created_at        time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	Updated_at        time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`

	Article []Article
}

type Reset_password struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	User_id   int       `json:"user_id" gorm:"not null"`
	Code      string    `json:"code" gorm:"type:varchar(255);not null"`
	Expired   time.Time `json:"expired" gorm:"type:datetime;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoCreateTime:milli"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;default:CURRENT_TIMESTAMP();autoUpdateTime:milli"`
}

type UserSignUp struct {
	Username        string `json:"username" validate:"required,min=3,max=50"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserSignIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResetPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotResetPasswordInput struct {
	NewPassword string `json:"new_password" validate:"required"`
}

type UserUpdateProfile struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
}
