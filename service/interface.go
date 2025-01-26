package service

import (
	"context"

	"github.com/JulianaSau/carzone/models"
)

type CarServiceInterface interface {
	GetCarById(ctx context.Context, id string) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, carReq *models.CarRequest) (*models.Car, error)
	UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (*models.Car, error)
	DeleteCar(ctx context.Context, id string) (*models.Car, error)
}

type EngineServiceInterface interface {
	GetEngineById(ctx context.Context, id string) (*models.Engine, error)
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (*models.Engine, error)
	UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (*models.Engine, error)
	DeleteEngine(ctx context.Context, id string) (*models.Engine, error)
}

type UserServiceInterface interface {
	GetUserProfile(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, userReq *models.UserRequest) (*models.User, error)
	UpdateUserProfile(ctx context.Context, id string, userReq *models.UserRequest) (*models.User, error)
	UpdateUserPassword(ctx context.Context, id string, userReq *models.UpdatePasswordRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id string) (*models.User, error)
	ToggleUserStatus(ctx context.Context, id string, active bool) (*models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type DriverServiceInterface interface {
	GetDrivers(ctx context.Context) ([]models.Driver, error)
	GetDriverById(ctx context.Context, id string) (*models.Driver, error)
	CreateDriver(ctx context.Context, driverReq *models.DriverRequest) (*models.Driver, error)
	UpdateDriver(ctx context.Context, id string, driverReq *models.DriverUpdateRequest) (*models.Driver, error)
	DeleteDriver(ctx context.Context, id string) (*models.Driver, error)
	SoftDeleteDriver(ctx context.Context, id string) (*models.Driver, error)
	ToggleDriverStatus(ctx context.Context, id string, active bool) (*models.Driver, error)
}
type TripServiceInterface interface {
	GetTrips(ctx context.Context) ([]models.Trip, error)
	GetTripsByDriverID(ctx context.Context, id string) ([]models.Trip, error)
	GetTripsByCarID(ctx context.Context, id string) ([]models.Trip, error)
	GetTripById(ctx context.Context, id string) (*models.Trip, error)
	CreateTrip(ctx context.Context, tripReq *models.TripRequest) (*models.Trip, error)
	UpdateTrip(ctx context.Context, id string, tripReq *models.TripRequest) (*models.Trip, error)
	UpdateTripStatus(ctx context.Context, id string, status string) (*models.Trip, error)
	DeleteTrip(ctx context.Context, id string) (*models.Trip, error)
}
