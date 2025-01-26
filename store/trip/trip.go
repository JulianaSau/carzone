package trip

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/JulianaSau/carzone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type TripStore struct {
	db *sql.DB
}

func New(db *sql.DB) *TripStore {
	return &TripStore{db: db}
}

func (u TripStore) GetTrips(ctx context.Context) ([]models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "GetTrips-Store")
	defer span.End()

	trips := []models.Trip{}

	query := `
		SELECT id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by
		FROM trip
	`
	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(
			&trip.ID,
			&trip.Description,
			&trip.DriverID,
			&trip.CarID,
			&trip.StartLocation,
			&trip.EndLocation,
			&trip.StartTime,
			&trip.EndTime,
			&trip.DistanceKM,
			&trip.FuelConsumedLiters,
			&trip.Status,
			&trip.CreatedAt,
			&trip.UpdatedAt,
			&trip.CreatedBy,
			&trip.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}
func (u TripStore) GetTripsByCarID(ctx context.Context, id string) ([]models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "GetTrips-Store")
	defer span.End()

	// Parse the car ID
	carID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid Car ID : %w", err)
	}

	trips := []models.Trip{}

	query := `
		SELECT id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by
		FROM trip
		WHERE car_id = $1
	`
	rows, err := u.db.QueryContext(ctx, query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(
			&trip.ID,
			&trip.Description,
			&trip.DriverID,
			&trip.CarID,
			&trip.StartLocation,
			&trip.EndLocation,
			&trip.StartTime,
			&trip.EndTime,
			&trip.DistanceKM,
			&trip.FuelConsumedLiters,
			&trip.Status,
			&trip.CreatedAt,
			&trip.UpdatedAt,
			&trip.CreatedBy,
			&trip.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}
func (u TripStore) GetTripsByDriverID(ctx context.Context, id string) ([]models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "GetTripsByDriverID-Store")
	defer span.End()

	// Parse the driver ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid Driver ID : %w", err)
	}

	trips := []models.Trip{}

	query := `
		SELECT id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by
		FROM trip
		WHERE driver_id = $1
	`
	rows, err := u.db.QueryContext(ctx, query, driverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var trip models.Trip
		err := rows.Scan(
			&trip.ID,
			&trip.Description,
			&trip.DriverID,
			&trip.CarID,
			&trip.StartLocation,
			&trip.EndLocation,
			&trip.StartTime,
			&trip.EndTime,
			&trip.DistanceKM,
			&trip.FuelConsumedLiters,
			&trip.Status,
			&trip.CreatedAt,
			&trip.UpdatedAt,
			&trip.CreatedBy,
			&trip.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

func (e *TripStore) GetTripById(ctx context.Context, id string) (models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "GetTripById-Store")
	defer span.End()

	var trip models.Trip

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return trip, nil
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(ctx, `SELECT id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by
	from trip 
	WHERE id=$1`,
		id).Scan(&trip.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return trip, nil
		}
		return trip, err
	}
	return trip, err

}

func (e *TripStore) CreateTrip(ctx context.Context, tripReq *models.TripRequest) (models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "CreateTrip-Store")
	defer span.End()

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Trip{}, nil
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	tripID := uuid.New()
	_, err = tx.ExecContext(ctx,
		`
		INSERT INTO trip (id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 , $10, $11, $12, $13, $14, $15)
	`, tripID,
		tripReq.Description,
		tripReq.DriverID,
		tripReq.CarID,
		tripReq.StartLocation,
		tripReq.EndLocation,
		tripReq.StartTime,
		tripReq.EndTime,
		tripReq.DistanceKM,
		tripReq.FuelConsumedLiters,
		tripReq.Status,
		time.Now(),
		time.Now(),
		"",
		"",
	)

	if err != nil {
		return models.Trip{}, err
	}

	trip := models.Trip{
		ID:                 tripID,
		Description:        tripReq.Description,
		DriverID:           tripReq.DriverID,
		CarID:              tripReq.CarID,
		StartLocation:      tripReq.StartLocation,
		EndLocation:        tripReq.EndLocation,
		StartTime:          tripReq.StartTime,
		EndTime:            tripReq.EndTime,
		DistanceKM:         tripReq.DistanceKM,
		FuelConsumedLiters: tripReq.FuelConsumedLiters,
		Status:             tripReq.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          "",
		UpdatedBy:          "",
	}

	return trip, nil
}

func (e *TripStore) UpdateTrip(ctx context.Context, id string, tripReq *models.TripRequest) (models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "UpdateTrip-Store")
	defer span.End()

	// Parse the trip ID
	tripID, err := uuid.Parse(id)
	if err != nil {
		return models.Trip{}, fmt.Errorf("invalid Trip ID : %w", err)
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Trip{}, nil
	}

	// Defer the rollback or commit
	defer func() {
		// if we find any problem with the transaction, we rollback
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			// if everything is fine, we commit the transaction
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	// Update the trip
	results, err := tx.ExecContext(ctx,
		`
	    UPDATE trip SET description=$1, driver_id=$2, car_id=$3, start_location=$4, end_location=$5, start_time=$6, end_time=$7, distance_km=$8, fuel_consumed_liters=$9, status=$10
		WHERE id=$11
		`,
		tripReq.Description, tripReq.DriverID, tripReq.CarID, tripReq.StartLocation, tripReq.EndLocation, tripReq.StartTime, tripReq.EndTime, tripReq.DistanceKM, tripReq.FuelConsumedLiters, tripReq.Status, tripID)

	if err != nil {
		return models.Trip{}, err
	}

	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Trip{}, err
	}
	if rowsAffected == 0 {
		return models.Trip{}, errors.New("no rows updated")
	}

	// Return the updated trip
	trip := models.Trip{
		ID:                 tripID,
		Description:        tripReq.Description,
		DriverID:           tripReq.DriverID,
		CarID:              tripReq.CarID,
		StartLocation:      tripReq.StartLocation,
		EndLocation:        tripReq.EndLocation,
		StartTime:          tripReq.StartTime,
		EndTime:            tripReq.EndTime,
		DistanceKM:         tripReq.DistanceKM,
		FuelConsumedLiters: tripReq.FuelConsumedLiters,
		Status:             tripReq.Status,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		CreatedBy:          "",
		UpdatedBy:          "",
	}

	return trip, nil
}
func (e *TripStore) UpdateTripStatus(ctx context.Context, id string, status string) (models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "UpdateTripStatus-Store")
	defer span.End()

	// Parse the trip ID
	tripID, err := uuid.Parse(id)
	if err != nil {
		return models.Trip{}, fmt.Errorf("invalid Trip ID : %w", err)
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Trip{}, nil
	}

	// Defer the rollback or commit
	defer func() {
		// if we find any problem with the transaction, we rollback
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			// if everything is fine, we commit the transaction
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	// Update the trip
	results, err := tx.ExecContext(ctx,
		`
	    UPDATE trip SET status=$1
		WHERE id=$2
		`,
		status, tripID)

	if err != nil {
		return models.Trip{}, err
	}

	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Trip{}, err
	}
	if rowsAffected == 0 {
		return models.Trip{}, errors.New("no rows updated")
	}

	// Return the updated trip
	trip := models.Trip{
		ID:        tripID,
		Status:    status,
		UpdatedAt: time.Now(),
	}

	return trip, nil
}

func (s *TripStore) DeleteTrip(ctx context.Context, id string) (models.Trip, error) {
	tracer := otel.Tracer("TripStore")
	ctx, span := tracer.Start(ctx, "DeleteTrip-Store")
	defer span.End()

	var trip models.Trip

	// Parse the trip ID
	tripID, err := uuid.Parse(id)
	if err != nil {
		return models.Trip{}, fmt.Errorf("invalid Trip ID: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Trip{}, nil
	}

	// Defer the rollback or commit
	defer func() {
		// if we find any problem with the transaction, we rollback
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			// if everything is fine, we commit the transaction
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	// check if the trip exists
	err = tx.QueryRowContext(ctx, `SELECT id, description, driver_id, car_id, start_location, end_location, start_time, end_time, distance_km, fuel_consumed_liters, status, created_at, updated_at, created_by, updated_by
	from trip 
	WHERE id=$1`,
		id).Scan(
		&trip.ID,
		&trip.Description,
		&trip.DriverID,
		&trip.CarID,
		&trip.StartLocation,
		&trip.EndLocation,
		&trip.StartTime,
		&trip.EndTime,
		&trip.DistanceKM,
		&trip.FuelConsumedLiters,
		&trip.Status,
		&trip.CreatedAt,
		&trip.UpdatedAt,
		&trip.CreatedBy,
		&trip.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return trip, nil
		}
		return trip, err
	}

	// Delete the trip
	result, err := tx.ExecContext(ctx,
		`DELETE FROM trip WHERE id=$1`, tripID)

	if err != nil {
		return models.Trip{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Trip{}, err
	}
	if rowsAffected == 0 {
		return models.Trip{}, errors.New("no rows were deleted")
	}

	return trip, nil
}
