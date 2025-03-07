package handlers

import (
	"Assignment1/consts"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Status represents the status response for the API
type Status struct {
	CountriesNowAPI  string `json:"CountriesNowStatus"`
	CountriesRestAPI string `json:"CountriesRestStatus"`
	Version          string `json:"Version"`
	Uptime           int64  `json:"Uptime"`
	Error            string `json:"error,omitempty"` // Only included if there's an error
}

// Variable to store the start time of the service for uptime calculation
var serviceStartTime int64

// InitializeUptime sets the start time of the service for uptime calculation
func InitializeUptime() {
	serviceStartTime = time.Now().Unix() // Set the service start time to the current Unix timestamp
}

// checkAPIStatus checks the status of an external API via GET or POST request
// Returns the status of the response or an error
func checkAPIStatus(url string, method string, payload *strings.Reader) (string, error) {
	var resp *http.Response
	var err error

	// Choose the appropriate HTTP method for the request
	if method == "POST" {
		resp, err = http.Post(url, "application/json", payload)
	} else {
		resp, err = http.Get(url) // Send a GET request
	}

	if err != nil {
		return "", err // Return both status and error
	}
	defer func(Body io.ReadCloser) {
		// Close the response body and handle any error that occurs
		err := Body.Close()
		if err != nil {
			// Log the error if closing the body fails
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	return resp.Status, nil
}

// StatusHandler handles requests to check the status of external APIs
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Only allows GET methods
	if r.Method != http.MethodGet {
		http.Error(w, "Status: "+r.Method+" method is not allowed. Use "+http.MethodGet+" method instead.", http.StatusMethodNotAllowed)
		return
	}

	// Check CountriesNow API status
	countriesNowStatus, errNow := checkAPIStatus(consts.COUNTRYNOWENDPOINT+"countries", "GET", nil)
	if errNow != nil {
		// If there's an error, store the error message in the status
		countriesNowStatus = "Error: " + errNow.Error()
	}

	// Check RESTCountries API status
	restCountriesStatus, errRest := checkAPIStatus(consts.COUNTRYRESTENDPOINT+"all", "GET", nil)
	if errRest != nil {
		// If there's an error, store the error message in the status
		restCountriesStatus = "Error: " + errRest.Error()
	}

	// Create status response
	status := Status{
		CountriesNowAPI:  countriesNowStatus,
		CountriesRestAPI: restCountriesStatus,
		Version:          "v1",
		Uptime:           time.Now().Unix() - serviceStartTime,
	}

	// Convert status to JSON
	jsonStatus, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		// If there's an error converting to JSON, return an internal server error
		http.Error(w, "Status: Error generating JSON response", http.StatusInternalServerError)
		return
	}

	// Set the content type of the response to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON status response to the response writer
	w.Write(jsonStatus)
}
