package car

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) *CarHandler {
	return &CarHandler{
		service: service,
	}
}

// GetCarByIDHandler godoc
// @Summary Get car by ID
// @Description Get a car by its ID
// @Tags Car
// @Accept  json
// @Produce  json
// @Param id path string true "Car ID"
// @Success 200 {object} models.Car
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Car not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cars/{id} [get]
// @Security Bearer
func (h *CarHandler) GetCarById(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarById-Handler")
	defer span.End()
	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// get the car by id from the car service
	res, err := h.service.GetCarById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting car by id: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car by id response: ", err)
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

// GetCarByBrandHandler godoc
// @Summary Get cars by brand
// @Description Get cars by brand
// @Tags Car
// @Accept  json
// @Produce  json
// @Param brand query string true "Car Brand"
// @Param isEngine query boolean false "Car with Engine"
// @Success 200 {object} []models.Car
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Car not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cars [get]
// @Security Bearer
func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "GetCarByBrand-Handler")
	defer span.End()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	resp, err := h.service.GetCarByBrand(ctx, brand, isEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting car by brand: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling car by brand response: ", err)
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

// CreateCarHandler godoc
// @Summary Create a new car
// @Description Create a new car
// @Tags Car
// @Accept  json
// @Produce  json
// @Param car body models.CarRequest true "Car Request"
// @Success 201 {object} models.Car
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cars [post]
// @Security Bearer
func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "CreateCar-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		log.Println("Error unmarshalling car request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the car
	createdCar, err := h.service.CreateCar(ctx, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating car: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(createdCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling created car response: ", err)
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

// UpdateCarHandler godoc
// @Summary Update a car
// @Description Update a car
// @Tags Car
// @Accept  json
// @Produce  json
// @Param id path string true "Car ID"
// @Param car body models.CarRequest true "Car Request"
// @Success 200 {object} models.Car
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cars/{id} [put]
// @Security Bearer
func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateCar-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var carReq models.CarRequest
	err = json.Unmarshal(body, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error unmarshalling car request: ", err)
		return
	}

	// update the car
	updatedCar, err := h.service.UpdateCar(ctx, id, &carReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating car: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated car response body: ", err)
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

// DeleteCarHandler godoc
// @Summary Delete a car
// @Description Delete a car
// @Tags Car
// @Accept  json
// @Produce  json
// @Param id path string true "Car ID"
// @Success 200 {object} models.Car
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/cars/{id} [delete]
// @Security Bearer
func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteCar-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	// delete the car
	deletedCar, err := h.service.DeleteCar(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting car: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(deletedCar)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling deleted car response body: ", err)
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
