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

func validateEngineRequest(engineReq EngineRequest) error {
	if err := validateDisplacement(engineReq.Displacement); err != nil {
		return err
	}
	if err := validateNoOfCylinders(engineReq.NoOfCylinders); err != nil {
		return err
	}
	if err := validateCarRanges(engineReq.CarRange); err != nil {
		return err
	}
	return nil
}

func validateDisplacement(displacement int64) error {
	if displacement <= 0 {
		return errors.New("displacement must be greater than 0")
	}
	return nil
}

func validateNoOfCylinders(noOfCylinders int64) error {
	if noOfCylinders <= 0 {
		return errors.New("number of cylinders must be greater than 0")
	}
	return nil
}

func validateCarRanges(carRange int64) error {
	if carRange <= 0 {
		return errors.New("car range must be greater than 0")
	}
	return nil
}
