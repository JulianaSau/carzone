package trip

import (
	"context"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/store"
	"go.opentelemetry.io/otel"
)

type TripService struct {
	store store.TripStoreInterface
}

func NewTripService(store store.TripStoreInterface) *TripService {
	return &TripService{
		store: store,
	}
}

func (s *TripService) GetTrips(ctx context.Context) ([]models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "GetTrips-Service")
	defer span.End()

	trips, err := s.store.GetTrips(ctx)
	if err != nil {
		return nil, err
	}
	return trips, nil
}
func (s *TripService) GetTripsByCarID(ctx context.Context, id string) ([]models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "GetTrips-Service")
	defer span.End()

	trips, err := s.store.GetTripsByCarID(ctx, id)
	if err != nil {
		return nil, err
	}
	return trips, nil
}
func (s *TripService) GetTripsByDriverID(ctx context.Context, id string) ([]models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "GetTrips-Service")
	defer span.End()

	trips, err := s.store.GetTripsByDriverID(ctx, id)
	if err != nil {
		return nil, err
	}
	return trips, nil
}
func (s *TripService) GetTripById(ctx context.Context, id string) (*models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "GetTripById-Service")
	defer span.End()

	trip, err := s.store.GetTripById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (s *TripService) CreateTrip(ctx context.Context, tripReq *models.TripRequest) (*models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "CreateTrip-Service")
	defer span.End()

	if err := models.ValidateTripRequest(*tripReq); err != nil {
		return nil, err
	}

	createdTrip, err := s.store.CreateTrip(ctx, tripReq)
	if err != nil {
		return nil, err
	}
	return &createdTrip, nil
}

func (s *TripService) UpdateTrip(ctx context.Context, id string, tripReq *models.TripRequest) (*models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "UpdateTrip-Service")
	defer span.End()

	if err := models.ValidateTripRequest(*tripReq); err != nil {
		return nil, err
	}

	updatedTrip, err := s.store.UpdateTrip(ctx, id, tripReq)
	if err != nil {
		return nil, err
	}
	return &updatedTrip, nil
}

func (s *TripService) UpdateTripStatus(ctx context.Context, id string, status string) (*models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "UpdateTrip-Service")
	defer span.End()

	updatedTrip, err := s.store.UpdateTripStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return &updatedTrip, nil
}

func (s *TripService) DeleteTrip(ctx context.Context, id string) (*models.Trip, error) {
	tracer := otel.Tracer("TripService")
	ctx, span := tracer.Start(ctx, "DeleteTrip-Service")
	defer span.End()

	deletedTrip, err := s.store.DeleteTrip(ctx, id)
	if err != nil {
		return nil, err
	}
	return &deletedTrip, nil
}
