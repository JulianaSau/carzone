package trip

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

type TripHandler struct {
	service service.TripServiceInterface
}

func NewTripHandler(service service.TripServiceInterface) *TripHandler {
	return &TripHandler{
		service: service,
	}
}

// GetTripsHandler godoc
// @Summary Get all trips
// @Description Get all trips
// @Tags Trip
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Trip
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips [get]
// @Security Bearer
func (h *TripHandler) GetTrips(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "GetTrips-Handler")
	defer span.End()

	trips, err := h.service.GetTrips(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting trips: ", err)
		return
	}

	body, err := json.Marshal(trips)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling trips response: ", err)
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

// GetTripByIDHandler godoc
// @Summary Get trip by ID
// @Description Get a trip by its ID
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Trip ID"
// @Success 200 {object} models.Trip
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Trip not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips/{id} [get]
// @Security Bearer
func (h *TripHandler) GetTripById(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "GetTripById-Handler")
	defer span.End()
	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// get the trip by id from the trip service
	res, err := h.service.GetTripById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting trip by id: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling trip by id response: ", err)
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

// GetTripByCarIDHandler godoc
// @Summary Get trips by Car ID
// @Description Get trips by car ID
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Trip ID"
// @Success 200 {object} models.Trip
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Trip not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/car/{id}/trips [get]
// @Security Bearer
func (h *TripHandler) GetTripsByCarID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "GetTripByCarId-Handler")
	defer span.End()
	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// get the trip by id from the trip service
	res, err := h.service.GetTripsByCarID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting trips by car id: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling trips by car id response: ", err)
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

// GetTripByDriverIDHandler godoc
// @Summary Get trips by Driver ID
// @Description Get trips by Driver ID
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} models.Trip
// @Failure 400 {string} string "Invalid Driver ID"
// @Failure 404 {string} string "Trips not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/driver/{id}/trips [get]
// @Security Bearer
func (h *TripHandler) GetTripsByDriverID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "GetTripByDriverId-Handler")
	defer span.End()
	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// get the trip by id from the trip service
	res, err := h.service.GetTripsByDriverID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting trips by driver id: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling trips by driver id response: ", err)
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

// CreateTripHandler godoc
// @Summary Create a new trip
// @Description Create a new trip
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param trip body models.TripRequest true "Trip Request"
// @Success 201 {object} models.Trip
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips [post]
// @Security Bearer
func (h *TripHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "CreateTrip-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var tripReq models.TripRequest
	err = json.Unmarshal(body, &tripReq)
	if err != nil {
		log.Println("Error unmarshalling trip request: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create the trip
	createdTrip, err := h.service.CreateTrip(ctx, &tripReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating trip: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(createdTrip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling created trip response: ", err)
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

// UpdateTripHandler godoc
// @Summary Update a trip
// @Description Update a trip
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Trip ID"
// @Param trip body models.TripRequest true "Trip Request"
// @Success 200 {object} models.Trip
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips/{id} [put]
// @Security Bearer
func (h *TripHandler) UpdateTrip(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateTrip-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var tripReq models.TripRequest
	err = json.Unmarshal(body, &tripReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error unmarshalling trip request: ", err)
		return
	}

	// update the trip
	updatedTrip, err := h.service.UpdateTrip(ctx, id, &tripReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating trip: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedTrip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated trip response body: ", err)
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

// DeleteTripHandler godoc
// @Summary Delete a trip
// @Description Delete a trip
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Trip ID"
// @Success 200 {object} models.Trip
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips/{id} [delete]
// @Security Bearer
func (h *TripHandler) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteTrip-Handler")
	defer span.End()

	params := mux.Vars(r)
	id := params["id"]

	// delete the trip
	deletedTrip, err := h.service.DeleteTrip(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting trip: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(deletedTrip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling deleted trip response body: ", err)
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

// ToggleTripStatusHandler godoc
// @Summary Toggle trip status
// @Description Toggle trip status by ID
// @Tags Trip
// @Accept  json
// @Produce  json
// @Param id path string true "Trip ID"
// @Param active query boolean true "Active status"
// @Success 200 {object} models.Trip
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Trip not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/v1/trips/{id}/update-status [put]
// @Security Bearer
func (h *TripHandler) UpdateTripStatus(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("TripHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateTripStatus-Handler")
	defer span.End()

	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]
	status := r.URL.Query().Get("status")

	// toggle the trip status
	toggledTrip, err := h.service.UpdateTripStatus(ctx, id, status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating trip status: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(toggledTrip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated trip response: ", err)
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
