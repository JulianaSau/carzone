package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrMissingField = errors.New("missing required field")

type Trip struct {
	ID                 uuid.UUID `json:"id"`                   // Unique trip identifier
	Description        string    `json:"description"`          //
	DriverID           uuid.UUID `json:"driver_id"`            // Reference to the Driver model
	CarID              uuid.UUID `json:"car_id"`               // Reference to the Car model
	StartLocation      string    `json:"start_location"`       // Starting point of the trip
	EndLocation        string    `json:"end_location"`         // Destination of the trip
	StartTime          time.Time `json:"start_time"`           // Trip start time
	EndTime            time.Time `json:"end_time"`             // Trip end time (nullable if still ongoing)
	DistanceKM         float64   `json:"distance_km"`          // Distance covered in kilometers
	FuelConsumedLiters float64   `json:"fuel_consumed_liters"` // Fuel consumed in liters
	Status             string    `json:"status"`               // Trip status (e.g., Completed, In Progress, Cancelled, Draft, Scheduled)
	CreatedAt          time.Time `json:"created_at"`           // Record creation timestamp
	UpdatedAt          time.Time `json:"updated_at"`           // Record last update timestamp
	CreatedBy          string    `json:"created_by"`           // User who created the record
	UpdatedBy          string    `json:"updated_by"`           // User who last updated the record
}

type TripRequest struct {
	Description        string    `json:"description"`
	DriverID           uuid.UUID `json:"driver_id"`
	CarID              uuid.UUID `json:"car_id"`
	StartLocation      string    `json:"start_location"`
	EndLocation        string    `json:"end_location"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	DistanceKM         float64   `json:"distance_km"`
	FuelConsumedLiters float64   `json:"fuel_consumed_liters"`
	Status             string    `json:"status"`
}

func ValidateTripRequest(tripReq TripRequest) error {
	if tripReq.Description == "" {
		return ErrMissingField
	}
	if tripReq.DriverID == uuid.Nil {
		return ErrMissingField
	}
	if tripReq.CarID == uuid.Nil {
		return ErrMissingField
	}
	if tripReq.StartLocation == "" {
		return ErrMissingField
	}
	if tripReq.EndLocation == "" {
		return ErrMissingField
	}
	if tripReq.StartTime.IsZero() {
		return ErrMissingField
	}
	if tripReq.DistanceKM == 0 {
		return ErrMissingField
	}
	if tripReq.FuelConsumedLiters == 0 {
		return ErrMissingField
	}
	if err := validateTripStatus(tripReq.Status); err != nil {
		return err
	}
	return nil
}

func validateTripStatus(status string) error {
	validateTripTypes := []string{"Completed", "Scheduled", "Ongoing", "Cancelled", "Draft"}
	for _, validType := range validateTripTypes {
		if status == validType {
			return nil
		}
	}
	return errors.New("status type must be one of: Completed, Scheduled, Ongoing, Cancelled, Draft")
}
