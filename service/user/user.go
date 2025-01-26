package user

import (
	"context"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/store"
	"go.opentelemetry.io/otel"
)

// holds reference to the user interface
// dependency injection - abstraction of data store operations for the user interface
type UserService struct {
	store store.UserStoreInterface
}

func NewUserService(store store.UserStoreInterface) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "GetUsers-Service")
	defer span.End()

	users, err := s.store.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, id string) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "GetUserById-Service")
	defer span.End()

	user, err := s.store.GetUserProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, userReq *models.UserRequest) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "CreateUser-Service")
	defer span.End()

	// if err := models.ValidateRequest(*userReq); err != nil {
	// 	return nil, err
	// }

	createdUser, err := s.store.CreateUser(ctx, userReq)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, id string, userReq *models.UserRequest) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "UpdateUserProfile-Service")
	defer span.End()
	// if err := models.ValidateRequest(*userReq); err != nil {
	// 	return nil, err
	// }

	updatedUser, err := s.store.UpdateUserProfile(ctx, id, userReq)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}
func (s *UserService) UpdateUserPassword(ctx context.Context, id string, userReq *models.UpdatePasswordRequest) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "UpdateUserPassword-Service")
	defer span.End()
	// if err := models.ValidateRequest(*userReq); err != nil {
	// 	return nil, err
	// }

	updatedUser, err := s.store.UpdateUserPassword(ctx, id, userReq)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (s *UserService) ToggleUserStatus(ctx context.Context, id string, active bool) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "ToggleUserStatus-Service")
	defer span.End()

	deletedUser, err := s.store.ToggleUserStatus(ctx, id, active)
	if err != nil {
		return nil, err
	}
	return &deletedUser, nil
}
func (s *UserService) DeleteUser(ctx context.Context, id string) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "DeleteUser-Service")
	defer span.End()

	deletedUser, err := s.store.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &deletedUser, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	tracer := otel.Tracer("UserService")
	ctx, span := tracer.Start(ctx, "GetUserByUsername-Service")
	defer span.End()

	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
