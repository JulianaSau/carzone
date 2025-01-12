package car

import (
	"context"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/store"
)

// holds reference to the car interface
// dependency injection - abstraction of data store operations for the car interface
type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{
		store: store,
	}
}

func (s *CarService) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	car, err := s.store.GetCarById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &car, nil
}

func (s *CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (s *CarService) CreateCar(ctx context.Context, carReq *models.CarRequest) (*models.Car, error) {
	if err := models.ValidateRequest(*carReq); err != nil {
		return nil, err
	}

	createdCar, err := s.store.CreateCar(ctx, carReq)
	if err != nil {
		return nil, err
	}
	return &createdCar, nil
}

func (s *CarService) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (*models.Car, error) {
	if err := models.ValidateRequest(*carReq); err != nil {
		return nil, err
	}

	updatedCar, err := s.store.UpdateCar(ctx, id, carReq)
	if err != nil {
		return nil, err
	}
	return &updatedCar, nil
}

func (s *CarService) DeleteCar(ctx context.Context, id string) (*models.Car, error) {
	deletedCar, err := s.store.DeleteCar(ctx, id)
	if err != nil {
		return nil, err
	}
	return &deletedCar, nil
}
