package store

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

type DriverStore struct {
	db *sql.DB
}

func New(db *sql.DB) *DriverStore {
	return &DriverStore{db: db}
}

func (d DriverStore) GetDrivers(ctx context.Context) ([]models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "GetDrivers-Store")
	defer span.End()

	drivers := []models.Driver{}

	query := `
		SELECT 
			d.id, d.user_id, d.driver_license_number, d.license_expiry, d.active,
			d.created_at, d.updated_at, d.created_by, d.deleted_at,
			u.id AS user_id, u.username, u.first_name, u.last_name, u.email
		FROM driver d
		JOIN user u ON d.user_id = u.id
		WHERE d.deleted_at IS NULL
	`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var driver models.Driver
		err := rows.Scan(
			&driver.ID,
			&driver.UserID,
			&driver.DriverLicenseNo,
			&driver.LicenseExpiry,
			&driver.Active,
			&driver.CreatedAt,
			&driver.UpdatedAt,
			&driver.CreatedBy,
			&driver.DeletedAt,
			&driver.User.ID,
			&driver.User.UserName,
			&driver.User.FirstName,
			&driver.User.LastName,
			&driver.User.Email,
		)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}
	return drivers, nil
}

func (d DriverStore) CreateDriver(ctx context.Context, driverReq *models.DriverRequest) (models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "CreateDriver-Store")
	defer span.End()
	var userId uuid.UUID
	// get user id from the user table
	err := d.db.QueryRowContext(
		ctx,
		`
				SELECT id 
				FROM user 
				WHERE id=$1
			`,
		driverReq.UserID,
	).Scan(&userId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Driver{}, errors.New("user_id doesnt exist in the user table")
		}
		return models.Driver{}, err
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Driver{}, nil
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

	// Insert driver into the database
	query := `
		INSERT INTO driver (id, user_id, driver_license_no, licence_expiry, active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid
	`
	driverID := uuid.New()
	_, err = tx.ExecContext(ctx, query,
		driverID,
		userId,
		driverReq.DriverLicenseNo,
		driverReq.LicenseExpiry,
		true,
		"d3b07384-d9a1-4c4b-8a0d-4b1b1b1b1b1b",
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return models.Driver{}, err
	}
	driver := models.Driver{
		ID:              driverID,
		UserID:          userId,
		DriverLicenseNo: driverReq.DriverLicenseNo,
		LicenseExpiry:   driverReq.LicenseExpiry,
		Active:          true,
		CreatedBy:       "d3b07384-d9a1-4c4b-8a0d-4b1b1b1b1b1b",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return driver, nil
}

func (d DriverStore) UpdateDriver(ctx context.Context, id string, driverReq *models.DriverUpdateRequest) (models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "UpdateDriver-Store")
	defer span.End()

	// Parse the driver ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return models.Driver{}, fmt.Errorf("invalid Driver ID : %w", err)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Driver{}, nil
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

	// Update driver profile in the database
	query := `
		UPDATE driver
		SET licence_expiry=$1, driver_license_no=$2, updated_at = $3
		WHERE uuid = $4
	`
	results, err := tx.ExecContext(ctx, query,
		driverReq.LicenseExpiry,
		driverReq.DriverLicenseNo,
		time.Now(),
		driverID,
	)

	if err != nil {
		return models.Driver{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Driver{}, err
	}
	if rowsAffected == 0 {
		return models.Driver{}, errors.New("no rows updated")
	}

	// Return the updated driver
	driver := models.Driver{
		ID:              driverID,
		DriverLicenseNo: driverReq.DriverLicenseNo,
		LicenseExpiry:   driverReq.LicenseExpiry,
		UpdatedAt:       time.Now(),
	}

	return driver, nil
}

func (d DriverStore) GetDriverById(ctx context.Context, id string) (models.Driver, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "GetDriverById-Store")
	defer span.End()

	// Parse the driver ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return models.Driver{}, fmt.Errorf("invalid driver ID : %w", err)
	}

	// Query the database
	query := `
	SELECT 
			d.id, d.user_id, d.driver_license_number, d.license_expiry, d.active,
			d.created_at, d.updated_at, d.created_by, d.deleted_at,
			u.id AS user_id, u.username, u.first_name, u.last_name, u.email
		FROM driver d
		JOIN user u ON d.user_id = u.id
		WHERE d.deleted_at IS NULL AND d.id = $1
	`
	row := d.db.QueryRowContext(ctx, query, driverID)

	driver := models.Driver{}
	err = row.Scan(
		&driver.ID,
		&driver.UserID,
		&driver.DriverLicenseNo,
		&driver.LicenseExpiry,
		&driver.Active,
		&driver.CreatedAt,
		&driver.UpdatedAt,
		&driver.CreatedBy,
		&driver.DeletedAt,
		&driver.User.ID,
		&driver.User.UserName,
		&driver.User.FirstName,
		&driver.User.LastName,
		&driver.User.Email,
	)
	if err != nil {
		return models.Driver{}, err
	}
	return driver, nil
}

func (d DriverStore) ToggleDriverStatus(ctx context.Context, id string, active bool) (models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "ToggleDriverStatus-Store")
	defer span.End()

	// Parse the user ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return models.Driver{}, fmt.Errorf("invalid driver ID : %w", err)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Driver{}, nil
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

	query := `
	    UPDATE driver
		SET active = $1, updated_at = $2
		WHERE uuid = $3
	`
	results, err := tx.ExecContext(ctx, query,
		active,
		time.Now(),
		driverID,
	)
	if err != nil {
		return models.Driver{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Driver{}, err
	}
	if rowsAffected == 0 {
		return models.Driver{}, errors.New("no rows updated")
	}
	// Return the updated user
	driver := models.Driver{
		ID:        driverID,
		Active:    active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return driver, nil
}

func (d DriverStore) DeleteDriver(ctx context.Context, id string) (models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "DeleteDriver-Store")
	defer span.End()

	// Parse the user ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return models.Driver{}, fmt.Errorf("invalid Driver ID : %w", err)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Driver{}, nil
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

	// check if the user exists
	err = tx.QueryRowContext(ctx, `SELECT id FROM driver WHERE id = $1`, driverID).Scan()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Driver{}, nil
		}
		return models.Driver{}, err
	}

	query := `
		DELETE FROM driver
		WHERE id = $1`

	results, err := tx.ExecContext(ctx, query, driverID)
	if err != nil {
		return models.Driver{}, err
	}
	// Check if the row was actually deleted
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Driver{}, err
	}
	if rowsAffected == 0 {
		return models.Driver{}, errors.New("no rows deleted")
	}
	// Return the deleted driver
	driver := models.Driver{
		ID:        driverID,
		CreatedAt: time.Now(),
	}

	return driver, nil
}

func (d DriverStore) SoftDeleteDriver(ctx context.Context, id string) (models.Driver, error) {
	tracer := otel.Tracer("DriverStore")
	ctx, span := tracer.Start(ctx, "SoftDeleteDriver-Store")
	defer span.End()

	// Parse the user ID
	driverID, err := uuid.Parse(id)
	if err != nil {
		return models.Driver{}, fmt.Errorf("invalid driver ID : %w", err)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Driver{}, nil
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
	query := `
	    UPDATE driver
		SET deleted_at = $1
		WHERE uuid = $2
	`
	results, err := tx.ExecContext(ctx, query,
		time.Now(),
		driverID,
	)
	if err != nil {
		return models.Driver{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Driver{}, err
	}
	if rowsAffected == 0 {
		return models.Driver{}, errors.New("no rows updated")
	}
	// Return the updated driver
	driver := models.Driver{
		ID:        driverID,
		DeletedAt: time.Now(),
	}
	return driver, nil
}
