package handlers

import (
	"Assignment1/consts"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FinalWrapper lagrer filtrerte befolkningsdata og gjennomsnittlig befolkning
type FinalWrapper struct {
	Mean   int `json:"mean"`
	Values []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	} `json:"values"`
}

// IsoStruct brukes til å hente ISO3-koder fra API
type IsoStruct struct {
	Iso3 string `json:"cca3"`
}

// Wrapper lagrer befolkningsdata fra API
type Wrapper struct {
	Mean   int `json:"mean"`
	Values struct {
		PopulationCounts []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// PopulationHandler håndterer forespørsler om befolkning
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	// Henter landkode fra query eller URL-path
	iso := r.URL.Query().Get("iso")
	if iso == "" {
		iso = strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/population/")
	}

	// Validerer at landkoden er på 2 bokstaver
	if len(iso) != 2 {
		http.Error(w, "Invalid country code: Must be 2 letters", http.StatusBadRequest)
		return
	}

	// Konverterer ISO2 til ISO3
	iso3 := ConvertIso(w, iso)
	if iso3 == "" {
		http.Error(w, "Invalid country code", http.StatusNotFound)
		return
	}

	// Henter tidsperiode (aar) fra query-parameter
	query := r.URL.Query().Get("limit")
	var start, end int
	if query == "" {
		start = 0
		end = time.Now().Year()
	} else {
		limit := strings.Split(query, "-")
		if len(limit) != 2 {
			http.Error(w, "Wrong input for query. You need 2 numbers.", http.StatusBadRequest)
			fmt.Fprintln(w, "expected 2 arguments in limit, got", len(limit))
			return
		}
		s, errConvStart := strconv.Atoi(limit[0])
		if errConvStart != nil {
			http.Error(w, "Start year must be an integer", http.StatusBadRequest)
			return
		} else {
			start = s
		}
		e, errConvEnd := strconv.Atoi(limit[1])
		if errConvEnd != nil {
			http.Error(w, "End year must be an integer", http.StatusBadRequest)
			return
		} else {
			end = e
		}
	}
	var popTemp FinalWrapper
	err := FetchPopulation(w, iso3, start, end, &popTemp)
	if err != nil {
		return
	}

	// Returnerer befolkningsdata som JSON
	jsonStatus, errjson := json.MarshalIndent(popTemp, "", "    ")
	if errjson != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(jsonStatus))
}

// ConvertIso konverterer en ISO2-landkode til ISO3 ved hjelp av et API
func ConvertIso(w http.ResponseWriter, iso string) string {
	resp, errGet := http.Get(consts.COUNTRYRESTENDPOINT + "alpha/" + iso + "?fields=cca3")
	if errGet != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Iso2 code is not in use.", http.StatusNotFound)
		return ""
	}

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return ""
	}
	var isoStructTemp IsoStruct
	errJson := json.Unmarshal(body, &isoStructTemp)
	if errJson != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return ""
	}

	if isoStructTemp.Iso3 == "" {
		http.Error(w, "Iso3 could not be retrieved from iso2 code \""+iso+"\".", http.StatusNotFound)
		return ""
	}

	return isoStructTemp.Iso3
}

// FetchPopulation henter befolkningsdata fra API og filtrerer dem etter årstall
func FetchPopulation(w http.ResponseWriter, iso3 string, min, max int, popTemp *FinalWrapper) error {

	payloadKEY := map[string]string{"iso3": iso3}
	payloadJSON, err := json.Marshal(payloadKEY)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return errors.New("Internal server error")
	}
	payload := strings.NewReader(string(payloadJSON))

	// Sender POST-request til API for befolkningsdata
	resp, err := http.Post(consts.COUNTRYNOWENDPOINT+"countries/population", "application/json", payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return errors.New("Internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Error: iso-2 code is not in use, pleace use iso-2.", http.StatusNotFound)
		return nil
	}

	body, errReadAll := io.ReadAll(resp.Body)
	if errReadAll != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return errors.New("Internal server error")
	}

	var wrapperTemp Wrapper
	errJson := json.Unmarshal(body, &wrapperTemp)

	if errJson != nil {
		log.Println("There was an error parsing json: ", errJson.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// Filtrerer befolkningsdata etter årstall
	var meanSum int
	// Går gjennom befolkningsdataene i wrapperTemp
	for i := 0; i < len(wrapperTemp.Values.PopulationCounts); i++ {
		if wrapperTemp.Values.PopulationCounts[i].Year >= min && wrapperTemp.Values.PopulationCounts[i].Year <= max {
			popTemp.Values = append(popTemp.Values, wrapperTemp.Values.PopulationCounts[i])
			meanSum += wrapperTemp.Values.PopulationCounts[i].Value
		}
	}
	// Beregner gjennomsnittlig befolkning
	if len(popTemp.Values) == 0 {
		popTemp.Mean = 0
	} else {
		popTemp.Mean = meanSum / len(popTemp.Values)
	}
	return nil
}
