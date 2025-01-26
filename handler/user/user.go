package user

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUsersHandler godoc
// @Summary Get all users
// @Description Get all users
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users [get]
// @Security Bearer
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "GetUsers-Handler")
	defer span.End()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting users: ", err)
		return
	}

	body, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling users response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

// GetUserProfileHandler godoc
// @Summary Get user profile
// @Description Get user profile by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id} [get]
// @Security Bearer
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "GetUserProfile-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUserProfile(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting user profile: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling user profile response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

// CreateUserHandler godoc
// @Summary Create user
// @Description Create a new user
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.UserRequest true "User object that needs to be created"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users [post]
// @Security Bearer
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "CreateUser-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var userReq models.UserRequest
	err = json.Unmarshal(body, &userReq)
	if err != nil {
		log.Println("Error unmarshalling user request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the car
	createdCar, err := h.service.CreateUser(ctx, &userReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating user: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling created user response: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// write the response body
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

// UpdateUserProfileHandler godoc
// @Summary Update user profile
// @Description Update user profile by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body models.UserRequest true "User object that needs to be updated"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id} [put]
// @Security Bearer
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateUserProfile-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var userReq models.UserRequest
	err = json.Unmarshal(body, &userReq)
	if err != nil {
		log.Println("Error unmarshalling user request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the user profile
	updatedUser, err := h.service.UpdateUserProfile(ctx, id, &userReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating user profile: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated user response: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("Error writing response body: ", err)
		return
	}
}

// UpdateUserPasswordHandler godoc
// @Summary Update user password
// @Description Update user password by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body models.UpdatePasswordRequest true "User object that needs to be updated"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id}/update-password [put]
// @Security Bearer
func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {

	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateUserPassword-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var userReq models.UpdatePasswordRequest
	err = json.Unmarshal(body, &userReq)
	if err != nil {
		log.Println("Error unmarshalling user request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the user profile
	updatedUser, err := h.service.UpdateUserPassword(ctx, id, &userReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating user profile: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated user response: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("Error writing response body: ", err)
		return
	}
}

// DeleteUserHandler godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id}/delete [delete]
// @Security Bearer
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteUser-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// delete the user
	deletedUser, err := h.service.DeleteUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting user: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(deletedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling deleted user response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

// ToggleUserStatusHandler godoc
// @Summary Toggle user status
// @Description Toggle user status by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param active query boolean true "Active status"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/users/{id}/toggle-status [put]
// @Security Bearer
func (h *UserHandler) ToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("UserHandler")
	ctx, span := tracer.Start(r.Context(), "ToggleUserStatus-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]
	active := r.URL.Query().Get("active")

	// parse the active status
	isActive, err := strconv.ParseBool(active)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid active status: ", err)
		return
	}

	// toggle the user status
	toggledUser, err := h.service.ToggleUserStatus(ctx, id, isActive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error toggling user status: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(toggledUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling toggled user response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}
