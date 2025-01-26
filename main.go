package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/JulianaSau/carzone/driver"
	carHandler "github.com/JulianaSau/carzone/handler/car"
	driverHandler "github.com/JulianaSau/carzone/handler/driver"
	engineHandler "github.com/JulianaSau/carzone/handler/engine"
	tripHandler "github.com/JulianaSau/carzone/handler/trip"
	userHandler "github.com/JulianaSau/carzone/handler/user"
	carService "github.com/JulianaSau/carzone/service/car"
	driverService "github.com/JulianaSau/carzone/service/driver"
	engineService "github.com/JulianaSau/carzone/service/engine"
	tripService "github.com/JulianaSau/carzone/service/trip"
	userService "github.com/JulianaSau/carzone/service/user"
	carStore "github.com/JulianaSau/carzone/store/car"
	driverStore "github.com/JulianaSau/carzone/store/driver"
	engineStore "github.com/JulianaSau/carzone/store/engine"
	tripStore "github.com/JulianaSau/carzone/store/trip"
	userStore "github.com/JulianaSau/carzone/store/user"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	loginHandler "github.com/JulianaSau/carzone/handler/login"
	middleware "github.com/JulianaSau/carzone/middleware"

	_ "github.com/JulianaSau/carzone/docs" // Import generated Swagger docs
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Car Management System API
// @version 1.0
// @description API documentation for the car management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	fmt.Println("Hello, CarZone!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	// start tracing
	traceProvider, err := startTracing()
	if err != nil {
		log.Fatalf("failed to start tracing: %v", err)
	}

	// shutdown the trace provider
	defer func() {
		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown trace provider: %v", err)
		}
	}()

	otel.SetTracerProvider(traceProvider)

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()
	// if db == nil {
	// 	log.Fatal("Could not connect to the database")
	// }

	// create a new car store instance and a new car service instance using the db instance
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	userStore := userStore.New(db)
	userService := userService.NewUserService(userStore)

	driverStore := driverStore.New(db)
	driverService := driverService.NewDriverService(driverStore)

	tripStore := tripStore.New(db)
	tripService := tripService.NewTripService(tripStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)
	userHandler := userHandler.NewUserHandler(userService)
	driverHandler := driverHandler.NewDriverHandler(driverService)
	tripHandler := tripHandler.NewTripHandler(tripService)

	// initialise router
	router := mux.NewRouter()

	// // execute schema file to populate database
	// schemaFile := "store/schema.sql"
	// if err := executeSchemaFile(db, schemaFile); err != nil {
	// 	log.Fatalf("Error while executing schema file: %v", err)
	// }

	// define routes
	router.Use(otelmux.Middleware("carzone"))
	router.Use(middleware.MetricsMiddleware)

	router.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler.LoginHandler(w, r, userService)
	}).Methods("POST")

	// Swagger documentation route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// middleware
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.AuthMIddleware)
	// router.Use(middleware.AuthMIddleware)

	protected.HandleFunc("/api/v1/users", userHandler.GetUsers).Methods("GET")
	protected.HandleFunc("/api/v1/users/{id}", userHandler.GetUserProfile).Methods("GET")
	protected.HandleFunc("/api/v1/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/api/v1/users/{id}", userHandler.UpdateUserProfile).Methods("PUT")
	protected.HandleFunc("/api/v1/users/{id}/update-password", userHandler.UpdateUserPassword).Methods("PUT")
	protected.HandleFunc("/api/v1/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	protected.HandleFunc("/api/v1/users/{id}/toggle-status", userHandler.ToggleUserStatus).Methods("PUT")

	protected.HandleFunc("/api/v1/drivers", driverHandler.GetDrivers).Methods("GET")
	protected.HandleFunc("/api/v1/drivers/{id}", driverHandler.GetDriverById).Methods("GET")
	protected.HandleFunc("/api/v1/drivers", driverHandler.CreateDriver).Methods("POST")
	protected.HandleFunc("/api/v1/drivers/{id}", driverHandler.UpdateDriver).Methods("PUT")
	protected.HandleFunc("/api/v1/drivers/{id}/delete", driverHandler.DeleteDriver).Methods("DELETE")
	protected.HandleFunc("/api/v1/drivers/{id}", driverHandler.SoftDeleteDriver).Methods("DELETE")
	protected.HandleFunc("/api/v1/drivers/{id}/toggle-status", driverHandler.ToggleDriverStatus).Methods("PUT")

	protected.HandleFunc("/api/v1/cars/{id}", carHandler.GetCarById).Methods("GET")
	protected.HandleFunc("/api/v1/cars", carHandler.GetCarByBrand).Methods("GET")
	protected.HandleFunc("/api/v1/cars", carHandler.CreateCar).Methods("POST")
	protected.HandleFunc("/api/v1/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	protected.HandleFunc("/api/v1/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	protected.HandleFunc("/api/v1/engines/{id}", engineHandler.GetEngineById).Methods("GET")
	protected.HandleFunc("/api/v1/engines", engineHandler.CreateEngine).Methods("POST")
	protected.HandleFunc("/api/v1/engines/{id}", engineHandler.UpdateEngine).Methods("PUT")
	protected.HandleFunc("/api/v1/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	protected.HandleFunc("/api/v1/trips", tripHandler.GetTrips).Methods("GET")
	protected.HandleFunc("/api/v1/trips/{id}", tripHandler.GetTripById).Methods("GET")
	protected.HandleFunc("/api/v1/cars/{id}/trips", tripHandler.GetTripsByCarID).Methods("GET")
	protected.HandleFunc("/api/v1/drivers/{id}/trips", tripHandler.GetTripsByDriverID).Methods("GET")
	protected.HandleFunc("/api/v1/trips", tripHandler.CreateTrip).Methods("POST")
	protected.HandleFunc("/api/v1/trips/{id}", tripHandler.UpdateTrip).Methods("PUT")
	protected.HandleFunc("/api/v1/trips/{id}/update-status", tripHandler.UpdateTripStatus).Methods("PUT")
	protected.HandleFunc("/api/v1/trips/{id}", tripHandler.DeleteTrip).Methods("DELETE")

	// metrics
	router.Handle("/metrics", promhttp.Handler())

	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}

func executeSchemaFile(db *sql.DB, fileName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Use filepath for cross-platform compatibility
	schemaPath := filepath.Clean(fileName)

	// Read the schema file
	sqlFile, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file '%s': %w", schemaPath, err)
	}

	// Normalize line endings
	sqlFile = bytes.ReplaceAll(sqlFile, []byte("\r\n"), []byte("\n"))

	// Execute the SQL commands
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("failed to execute schema file '%s': %w", schemaPath, err)
	}

	log.Printf("Successfully executed schema file: %s", schemaPath)
	return nil
}

func startTracing() (*trace.TracerProvider, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	// http exporter that will send data to Jaeggar
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("jaeger:4318"),
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create new Jaeger exporter: %w", err)
	}

	// create a new trace provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("carzone"),
			),
		),
	)
	return tracerProvider, nil
}
