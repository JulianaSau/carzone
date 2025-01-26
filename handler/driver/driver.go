package driver

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

type DriverHandler struct {
	service service.DriverServiceInterface
}

func NewDriverHandler(service service.DriverServiceInterface) *DriverHandler {
	return &DriverHandler{
		service: service,
	}
}

// GetDriversHandler godoc
// @Summary Get all drivers
// @Description Get all drivers
// @Tags Driver
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Driver
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers [get]
// @Security Bearer
func (h *DriverHandler) GetDrivers(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "GetDrivers-Handler")
	defer span.End()

	drivers, err := h.service.GetDrivers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting drivers: ", err)
		return
	}

	body, err := json.Marshal(drivers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling drivers response: ", err)
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

// GetDriverProfileHandler godoc
// @Summary Get driver profile
// @Description Get driver profile by ID
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers/{id} [get]
// @Security Bearer
func (h *DriverHandler) GetDriverById(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "GetDriverById-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	driver, err := h.service.GetDriverById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting driver profile: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(driver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling driver profile response: ", err)
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

// CreateDriverHandler godoc
// @Summary Create driver
// @Description Create a new driver
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param driver body models.DriverRequest true "Driver object that needs to be created"
// @Success 201 {object} models.Driver
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers [post]
// @Security Bearer
func (h *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "CreateDriver-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var driverReq models.DriverRequest
	err = json.Unmarshal(body, &driverReq)
	if err != nil {
		log.Println("Error unmarshalling driver request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the car
	createdCar, err := h.service.CreateDriver(ctx, &driverReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating driver: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling created driver response: ", err)
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

// UpdateDriverProfileHandler godoc
// @Summary Update driver profile
// @Description Update driver profile by ID
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Param driver body models.DriverRequest true "Driver object that needs to be updated"
// @Success 200 {object} models.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers/{id} [put]
// @Security Bearer
func (h *DriverHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateDriver-Handler")
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

	var driverReq models.DriverUpdateRequest
	err = json.Unmarshal(body, &driverReq)
	if err != nil {
		log.Println("Error unmarshalling driver request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// update the driver profile
	updatedDriver, err := h.service.UpdateDriver(ctx, id, &driverReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating driver profile: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated driver response: ", err)
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

// DeleteDriverHandler godoc
// @Summary Delete driver
// @Description Delete driver by ID
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers/{id}/delete [delete]
// @Security Bearer
func (h *DriverHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteDriver-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// delete the driver
	deletedDriver, err := h.service.DeleteDriver(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting driver: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(deletedDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling deleted driver response: ", err)
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

// SoftDeleteDriverHandler godoc
// @Summary Soft Delete driver
// @Description Soft Delete driver by ID
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers/{id} [delete]
// @Security Bearer
func (h *DriverHandler) SoftDeleteDriver(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "SoftDeleteDriver-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// delete the driver
	deletedDriver, err := h.service.SoftDeleteDriver(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error soft deleting driver: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(deletedDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling soft deleted driver response: ", err)
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

// ToggleDriverStatusHandler godoc
// @Summary Toggle driver status
// @Description Toggle driver status by ID
// @Tags Driver
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Param active query boolean true "Active status"
// @Success 200 {object} models.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/drivers/{id}/toggle-status [put]
// @Security Bearer
func (h *DriverHandler) ToggleDriverStatus(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("DriverHandler")
	ctx, span := tracer.Start(r.Context(), "ToggleDriverStatus-Handler")
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

	// toggle the driver status
	toggledDriver, err := h.service.ToggleDriverStatus(ctx, id, isActive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error toggling driver status: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(toggledDriver)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling toggled driver response: ", err)
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
