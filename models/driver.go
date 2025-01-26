package models

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"` // Reference to the User model
	DriverLicenseNo string    `json:"driver_license_number"`
	LicenseExpiry   time.Time `json:"license_expiry"` // e.g., 2025-12-31T23:59:59Z
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       string    `json:"created_by"`
	DeletedAt       time.Time `json:"deleted_at"`
	// Embedding the User struct to access user information like name, email, etc.
	User User `json:"user"`
}

type DriverRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	DriverLicenseNo string    `json:"driver_license_number"`
	LicenseExpiry   time.Time `json:"license_expiry"`
}

type DriverUpdateRequest struct {
	DriverLicenseNo string    `json:"driver_license_number"`
	LicenseExpiry   time.Time `json:"license_expiry"`
}
