package main

import (
	"Assignment1/consts"
	"Assignment1/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialiserer uptime-monitoring
	handlers.InitializeUptime()

	// Definerer HTTP-ruter
	http.HandleFunc("/countryinfo/v1/population/", handlers.PopulationHandler)
	http.HandleFunc("/countryinfo/v1/info/", handlers.InfoHandler)
	http.HandleFunc("/countryinfo/v1/status", handlers.StatusHandler)

	// Starter serveren
	err := http.ListenAndServe(":"+consts.PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
