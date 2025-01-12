package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID      uuid.UUID `json:"engine_id"`
	Displacement  int64     `json:"displacement"`
	NoOfCylinders int64     `json:"no_of_cylinders"`
	CarRange      int64     `json:"car_range"`
}

type EngineRequest struct {
	Displacement  int64 `json:"displacement"`
	NoOfCylinders int64 `json:"no_of_cylinders"`
	CarRange      int64 `json:"car_range"`
}

func ValidateEngineRequest(engineReq EngineRequest) error {
	if err := ValidateDisplacement(engineReq.Displacement); err != nil {
		return err
	}
	if err := ValidateNoOfCylinders(engineReq.NoOfCylinders); err != nil {
		return err
	}
	if err := ValidateCarRanges(engineReq.CarRange); err != nil {
		return err
	}
	return nil
}

func ValidateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater than 0")
	}
	return nil
}

func ValidateNoOfCylinders(noOfCylinders int64) error {
	if noOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	return nil
}

func ValidateCarRanges(carRange int64) error {
	if carRange <= 0 {
		return errors.New("car range must be greater than 0")
	}
	return nil
}
