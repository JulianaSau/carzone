package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/JulianaSau/carzone/models"
	"github.com/JulianaSau/carzone/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (h *EngineHandler) GetEngineById(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "GetEngineById-Handler")
	defer span.End()
	// get the request params
	vars := mux.Vars(r)
	id := vars["id"]

	// get the engine by id from the engine service
	res, err := h.service.GetEngineById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting engine by id: ", err)
		return
	}

	// marshal the response
	body, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body'
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

func (h *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "CreateEngine-Handler")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var engineReq models.EngineRequest

	// unmarshal the request body
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling engine request: ", err)
		return
	}

	// create the engine
	createdEngine, err := h.service.CreateEngine(ctx, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error creating engine: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(createdEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling engine response: ", err)
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

func (h *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "UpdateEngine-Handler")
	defer span.End()

	// get the request params
	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error reading request body: ", err)
		return
	}

	var engineReq models.EngineRequest

	// unmarshal the request body
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error unmarshalling engine request: ", err)
		return
	}

	// update the engine
	updatedEngine, err := h.service.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating engine: ", err)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(updatedEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling updated engine response: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// write the response body
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}

func (h *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EngineHandler")
	ctx, span := tracer.Start(r.Context(), "DeleteEngine-Handler")
	defer span.End()

	// get the request params
	params := mux.Vars(r)
	id := params["id"]

	// delete the engine
	deletedEngine, err := h.service.DeleteEngine(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting engine: ", err)
		response := map[string]string{"error": "Invalid ID or engine not found"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	// check if engine was deleted properly
	if deletedEngine.EngineID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{"error": "Engine not found"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)
		return
	}

	// marshal the response
	responseBody, err := json.Marshal(deletedEngine)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling deleted engine response: ", err)
		response := map[string]string{"error": "Internal server error"}
		jsonResponse, _ := json.Marshal(response)
		_, _ = w.Write(jsonResponse)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("Error writing response body: ", err)
	}
}
