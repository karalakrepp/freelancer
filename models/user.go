package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            int                `json:"id"`
	FirstName     string             `json:"firstName"`
	LastName      string             `json:"lastName"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Password      string             `json:"password"`
	PhoneNumber   int                `json:"phoneNumber"`
	Location      CountryDataProfile `json:"location"`
	UserType      string             `json:"user_type"`
	IsAdmin       bool               `json:"is_admin"`
	IsActive      bool               `json:"is_active"`
	DeactivatedAt time.Time          `json:"deactivated_at"`
	IsDeleted     bool               `json:"is_deleted"`
	Balance       float64            `json:"balance"`
	DeletedAt     time.Time          `json:"deleted_at"`
	CreatedAt     time.Time          `json:"created_at"`
}

type UserProfile struct {
	ID     int       `json:"id"`
	Owner  OwnerInfo `json:"owner"`
	UserID int       `json:"user_id"`
	//hakkinda
	Description      string       `json:"description"`
	Title            string       `json:"title"`
	Skill            []UserSkills `json:"skill"`
	Picture          string       `json:"picture"`
	ProjectCompleted int          `json:"projectCompleted"`
}
type CreateUserRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `json:"username" validate:"required,min=7"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=7"`
	PhoneNumber int    `json:"phoneNumber"`
	Location    string `json:"location"`
}
type DeleteAccountRequest struct {
	UserID   int64  `json:"user_id" validate:"required,number,min=1"`
	Password string `json:"password" validate:"required,min=7"`
}
type DeleteAccountResponse struct {
	Message string `json:"message"`
}
type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}
type LoginUserRequest struct {
	UserID   int64  `json:"user_id" validate:"required,number,min=1"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`

	PhoneNumber int                `json:"phoneNumber"`
	Location    CountryDataProfile `json:"location"`
	UserType    string             `json:"user_type"`
	CreatedAt   time.Time          `json:"created_at"`
}
