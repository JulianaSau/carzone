package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JulianaSau/carzone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type EngineStore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineStore {
	return &EngineStore{db: db}
}

func (e *EngineStore) GetEngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "GetEngineById-Store")
	defer span.End()

	var engine models.Engine

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return engine, nil
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

	err = tx.QueryRowContext(ctx, `SELECT id, displacement, no_of_cylinders, car_range
	from engine 
	WHERE id=$1`,
		id).Scan(&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, nil
		}
		return engine, err
	}
	return engine, err

}

func (e *EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "CreateEngine-Store")
	defer span.End()

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, nil
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

	engineID := uuid.New()
	_, err = tx.ExecContext(ctx,
		`
		INSERT INTO engine (id, displacement, no_of_cylinders, car_range)
		VALUES ($1, $2, $3, $4)
	`, engineID, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange)

	if err != nil {
		return models.Engine{}, err
	}

	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil
}

func (e *EngineStore) UpdateEngine(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "UpdateEngine-Store")
	defer span.End()

	// Parse the engine ID
	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid Engine ID : %w", err)
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, nil
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

	// Update the engine
	results, err := tx.ExecContext(ctx,
		`
	    UPDATE engine SET displacement=$1, no_of_cylinders=$2, car_range=$3
		WHERE id=$4
		`,
		engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange, engineID)

	if err != nil {
		return models.Engine{}, err
	}

	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("no rows updated")
	}

	// Return the updated engine
	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil
}

func (s *EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "DeleteEngine-Store")
	defer span.End()

	var engine models.Engine

	// Parse the engine ID
	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid Engine ID: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, nil
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

	// check if the engine exists
	err = tx.QueryRowContext(ctx, `SELECT id, displacement, no_of_cylinders, car_range
	from engine 
	WHERE id=$1`,
		id).Scan(&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, nil
		}
		return engine, err
	}

	// Delete the engine
	result, err := tx.ExecContext(ctx,
		`DELETE FROM  engine WHERE id=$1`, engineID)

	if err != nil {
		return models.Engine{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("no rows were deleted")
	}

	return engine, nil
}
