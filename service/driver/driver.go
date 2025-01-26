package driver

import (
	"context"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/store"
	"go.opentelemetry.io/otel"
)

// holds reference to the driver interface
// dependency injection - abstraction of data store operations for the driver interface
type DriverService struct {
	store store.DriverStoreInterface
}

func NewDriverService(store store.DriverStoreInterface) *DriverService {
	return &DriverService{
		store: store,
	}
}

func (s *DriverService) GetDrivers(ctx context.Context) ([]models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "GetDrivers-Service")
	defer span.End()

	drivers, err := s.store.GetDrivers(ctx)
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

func (s *DriverService) GetDriverById(ctx context.Context, id string) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "GetDriverById-Service")
	defer span.End()

	driver, err := s.store.GetDriverById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &driver, nil
}

func (s *DriverService) CreateDriver(ctx context.Context, driverReq *models.DriverRequest) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "CreateDriver-Service")
	defer span.End()

	// if err := models.ValidateRequest(*driverReq); err != nil {
	// 	return nil, err
	// }

	createdDriver, err := s.store.CreateDriver(ctx, driverReq)
	if err != nil {
		return nil, err
	}
	return &createdDriver, nil
}

func (s *DriverService) UpdateDriver(ctx context.Context, id string, driverReq *models.DriverUpdateRequest) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "UpdateDriver-Service")
	defer span.End()
	// if err := models.ValidateRequest(*driverReq); err != nil {
	// 	return nil, err
	// }

	updatedDriver, err := s.store.UpdateDriver(ctx, id, driverReq)
	if err != nil {
		return nil, err
	}
	return &updatedDriver, nil
}

func (s *DriverService) ToggleDriverStatus(ctx context.Context, id string, active bool) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "ToggleDriverStatus-Service")
	defer span.End()

	deletedDriver, err := s.store.ToggleDriverStatus(ctx, id, active)
	if err != nil {
		return nil, err
	}
	return &deletedDriver, nil
}
func (s *DriverService) DeleteDriver(ctx context.Context, id string) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "DeleteDriver-Service")
	defer span.End()

	deletedDriver, err := s.store.DeleteDriver(ctx, id)
	if err != nil {
		return nil, err
	}
	return &deletedDriver, nil
}
func (s *DriverService) SoftDeleteDriver(ctx context.Context, id string) (*models.Driver, error) {
	tracer := otel.Tracer("DriverService")
	ctx, span := tracer.Start(ctx, "SoftDeleteDriver-Service")
	defer span.End()

	deletedDriver, err := s.store.SoftDeleteDriver(ctx, id)
	if err != nil {
		return nil, err
	}
	return &deletedDriver, nil
}
