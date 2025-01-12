package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/JulianaSau/carzone/driver"
	carHandler "github.com/JulianaSau/carzone/handler/car"
	engineHandler "github.com/JulianaSau/carzone/handler/engine"
	carService "github.com/JulianaSau/carzone/service/car"
	engineService "github.com/JulianaSau/carzone/service/engine"
	carStore "github.com/JulianaSau/carzone/store/car"
	engineStore "github.com/JulianaSau/carzone/store/engine"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Hello, CarZone!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()

	// create a new car store instance and a new car service instance using the db instance
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	// initialise router
	router := mux.NewRouter()

	// define routes
	router.HandleFunc("/api/v1/cars/{id}", carHandler.GetCarById).Methods("GET")
	router.HandleFunc("/api/v1/cars", carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/api/v1/cars", carHandler.CreateCar).Methods("POST")
	router.HandleFunc("/api/v1/cars/{id}", carHandler.UpdateCar).Methods("UPDATE")
	router.HandleFunc("/api/v1/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/api/v1/engines/{id}", engineHandler.GetEngineById).Methods("GET")
	router.HandleFunc("/api/v1/engines", engineHandler.CreateEngine).Methods("POST")
	router.HandleFunc("/api/v1/engines/{id}", engineHandler.UpdateEngine).Methods("UPDATE")
	router.HandleFunc("/api/v1/engines/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
