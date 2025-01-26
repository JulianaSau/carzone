package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName    string    `json:"username"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	ID          uuid.UUID `json:"uuid"`
	Active      bool      `json:"active"`
	CreatedBy   string    `json:"created_by"`
	DeletedAt   time.Time `json:"deleted_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRequest struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Role            string `json:"role"`
	// ID              uuid.UUID `json:"uuid"`
}
type UpdatePasswordRequest struct {
	PreviousPassword string `json:"previous_password"`
	Password         string `json:"password"`
	ConfirmPassword  string `json:"confirm_password"`
}

// create a function to hash the user's password as you are saving the user
func (user *UserRequest) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *UpdatePasswordRequest) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// create a function to check if the password is a match
// CheckPassword checks if the provided password matches the hashed password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
