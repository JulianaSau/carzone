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

type UserStore struct {
	db *sql.DB
}

func New(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (u UserStore) GetUsers(ctx context.Context) ([]models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "GetUsers-Store")
	defer span.End()

	users := []models.User{}

	query := `
		SELECT username, first_name, last_name, email, phone_number, role, id, active, created_at, updated_at
		FROM "user"
	`
	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PhoneNumber,
			&user.Role,
			&user.ID,
			&user.Active,
			// &user.CreatedBy,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u UserStore) CreateUser(ctx context.Context, userReq *models.UserRequest) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "CreateUser-Store")
	defer span.End()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, nil
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

	// compare password and compare password
	if userReq.Password != userReq.ConfirmPassword {
		return models.User{}, errors.New("password and confirm password do not match")
	}

	// Hash the user's password before saving

	if err := userReq.HashPassword(userReq.Password); err != nil {
		return models.User{}, err
	}

	// Insert user into the database
	query := `
		INSERT INTO "user" (id, username, password, first_name, last_name, email, phone_number, role, active, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, username, first_name, last_name, email, phone_number, role, active, created_at, updated_at
	`
	userID := uuid.New()
	_, err = tx.ExecContext(ctx, query,
		userID,
		userReq.UserName,
		userReq.Password,
		userReq.FirstName,
		userReq.LastName,
		userReq.Email,
		userReq.PhoneNumber,
		userReq.Role,
		true,
		"d3b07384-d9a1-4c4b-8a0d-4b1b1b1b1b1b",
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		UserName:    userReq.UserName,
		FirstName:   userReq.FirstName,
		LastName:    userReq.LastName,
		Email:       userReq.Email,
		PhoneNumber: userReq.PhoneNumber,
		Role:        userReq.Role,
		ID:          userID,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return user, nil
}

func (u UserStore) UpdateUserProfile(ctx context.Context, id string, userReq *models.UserRequest) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "UpdateUserProfile-Store")
	defer span.End()

	// Parse the user ID
	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, nil
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

	// Update user profile in the database
	query := `
		UPDATE "user"
		SET first_name = $1, last_name = $2, username = $3, email = $4, phone_number = $5, updated_at = $6
		WHERE id = $6
	`
	results, err := tx.ExecContext(ctx, query,
		userReq.FirstName,
		userReq.LastName,
		userReq.UserName,
		userReq.Email,
		userReq.PhoneNumber,
		time.Now(),
	)

	if err != nil {
		return models.User{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, errors.New("no rows updated")
	}

	// Return the updated user
	user := models.User{
		UserName:    userReq.UserName,
		FirstName:   userReq.FirstName,
		LastName:    userReq.LastName,
		Email:       userReq.Email,
		PhoneNumber: userReq.PhoneNumber,
		ID:          userID,
		Active:      true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Role:        userReq.Role,
	}

	return user, nil
}

func (u UserStore) UpdateUserPassword(ctx context.Context, id string, userReq *models.UpdatePasswordRequest) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "UpdateUserPassword-Store")
	defer span.End()

	// Parse the user ID
	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, nil
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
	// Hash the new password before saving
	if err := userReq.HashPassword(userReq.Password); err != nil {
		return models.User{}, err
	}

	// Update user password in the database
	query := `
		UPDATE "user"
		SET password = $1, updated_at = $2
		WHERE id = $3
	`
	results, err := tx.ExecContext(ctx, query,
		userReq.Password,
		time.Now(),
		userID,
	)
	if err != nil {
		return models.User{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, errors.New("no rows updated")
	}
	// Return the updated user
	user := models.User{
		ID:        userID,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return user, nil
}

func (u UserStore) GetUserProfile(ctx context.Context, id string) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "GetUserProfile-Store")
	defer span.End()

	// Parse the user ID
	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	}

	// Query the database
	query := `
		SELECT username, first_name, last_name, email, phone_number, role, id, active, created_by, created_at, updated_at
		FROM "user"
		WHERE id = $1
	`
	row := u.db.QueryRowContext(ctx, query, userID)

	user := models.User{}
	err = row.Scan(
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&user.Role,
		&user.ID,
		&user.Active,
		&user.CreatedBy,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (u UserStore) ToggleUserStatus(ctx context.Context, id string, active bool) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "ToggleUserStatus-Store")
	defer span.End()

	// Parse the user ID
	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, nil
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
	    UPDATE "user"
		SET active = $1, updated_at = $2
		WHERE id = $3
	`
	results, err := tx.ExecContext(ctx, query,
		active,
		time.Now(),
		userID,
	)
	if err != nil {
		return models.User{}, err
	}
	// Check if the row was actually updated
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, errors.New("no rows updated")
	}
	// Return the updated user
	user := models.User{
		ID:        userID,
		Active:    active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

func (u UserStore) DeleteUser(ctx context.Context, id string) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "DeleteUser-Store")
	defer span.End()

	// Parse the user ID
	userID, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return models.User{}, nil
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
	err = tx.QueryRowContext(ctx, `SELECT id, username FROM "user" WHERE id = $1`, userID).Scan()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, nil
		}
		return models.User{}, err
	}

	query := `
		DELETE FROM "user"
		WHERE id = $1`

	results, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return models.User{}, err
	}
	// Check if the row was actually deleted
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, errors.New("no rows deleted")
	}
	// Return the deleted user
	user := models.User{
		ID:        userID,
		CreatedAt: time.Now(),
	}

	return user, nil
}

func (u UserStore) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	tracer := otel.Tracer("UserStore")
	ctx, span := tracer.Start(ctx, "GetUserProfile-Store")
	defer span.End()

	// Parse the user name
	// _, err := uuid.Parse(username)
	// if err != nil {
	// 	return models.User{}, fmt.Errorf("invalid User ID : %w", err)
	// }

	// Prepare the SQL query to find a user by username
	query := `
		SELECT id, username, password, first_name, last_name, email, phone_number, role, active, created_at, updated_at
		FROM "user"
		WHERE username = $1 AND deleted_at IS NULL
	`
	row := u.db.QueryRowContext(ctx, query, username)

	user := models.User{}
	err := row.Scan(
		&user.ID,
		&user.UserName,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.User{}, nil
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
