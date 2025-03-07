package handlers

import (
	"Assignment1/consts"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Struktur for å holde landinformasjon
type Country struct {
	Name struct {
		Common string `json:"common"` // Vanlig navn på landet
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    []string          `json:"capital"`
	Cities     []string          `json:"data"`
}

// Struktur for å holde byene som vi henter fra API
type JutsCities struct {
	Cities []string `json:"data"`
}

// Håndterer forespørsler om landinformasjon
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	// Kun GET-metoden er tillatt
	if r.Method != http.MethodGet {
		http.Error(w, r.Method+" method is not allowed. Use "+http.MethodGet+" method instead.", http.StatusMethodNotAllowed)
		return
	}

	// Henter iso-koden fra URL-queryen
	iso := r.URL.Query().Get("iso")
	iso = strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/info/")
	if len(iso) != 2 {
		http.Error(w, "Iso-2 must be a 2 letter code", http.StatusBadRequest)
		return
	}

	// Hent landinformasjon ved hjelp av iso-koden
	country := Country{}
	FetchCountry(w, &country, iso)

	// Sjekk om vi fikk gyldig landinformasjon
	if country.Name.Common == "" {
		http.Error(w, "Could not fetch country data.", http.StatusInternalServerError)
		return
	}

	// Håndterer limit-parameteren for byene
	query := r.URL.Query().Get("limit")
	limit := 10 // Standardverdi for limit
	if query != "" {
		var err error
		limit, err = strconv.Atoi(query) // Konverterer limit til heltall
		if err != nil || limit <= 0 {
			http.Error(w, "Limit must be a positive integer.", http.StatusBadRequest)
			return
		}
	}

	// Hent byene med angitt limit
	FetchCities(w, &country, limit)

	// Skriv ut landinfo som JSON
	PrintCountry(w, country)
}

// FetchCountry henter landdata fra API ved hjelp av iso-koden
func FetchCountry(w http.ResponseWriter, c *Country, iso string) {
	// Hent landinfo fra API-endepunktet
	resp, err := http.Get(fmt.Sprintf("%salpha/%s", consts.COUNTRYRESTENDPOINT, iso))
	if err != nil {
		http.Error(w, "Error fetching country data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Sjekk om landet ikke finnes
	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Iso2 code is not in use.", http.StatusNotFound)
		return
	}

	// Les responsen fra API-et
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	// Parse JSON-responsen (det er et array)
	var countries []Country
	if err := json.Unmarshal(body, &countries); err != nil {
		log.Println("JSON Parsing Error:", err)
		http.Error(w, "Error parsing country data.", http.StatusInternalServerError)
		return
	}

	// Hvis vi fikk data, sett det første landet til country-objektet
	if len(countries) > 0 {
		*c = countries[0]
	} else {
		http.Error(w, "No country found for the provided iso code.", http.StatusNotFound)
		return
	}
}

// FetchCities henter byer fra CountriesNow API og filtrerer dem
func FetchCities(w http.ResponseWriter, c *Country, limit int) {
	// Sjekk om landets navn er tomt
	if c.Name.Common == "" {
		http.Error(w, "Country name is empty.", http.StatusInternalServerError)
		return
	}

	// Lag payload for POST-request
	payload := strings.NewReader(fmt.Sprintf("{\"country\": \"%s\"}", c.Name.Common))

	// Gjør en POST-request til CountriesNow API for å hente byer
	resp, err := http.Post(consts.COUNTRYNOWENDPOINT+"countries/cities", "application/json", payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Sjekk om landet ikke finnes
	if resp.StatusCode == http.StatusNotFound {
		http.Error(w, "Iso2 code is not in use.", http.StatusNotFound)
		return
	}

	// Les og parse JSON-responsen
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Struktur for å holde byene
	var temp JutsCities

	// Parse JSON-responsen
	if err := json.Unmarshal(body, &temp); err != nil {
		http.Error(w, "Error parsing cities data.", http.StatusInternalServerError)
		return
	}

	// Hvis ingen byer ble funnet, returner feil
	if len(temp.Cities) == 0 {
		http.Error(w, "There are no cities found in this country.", http.StatusNotFound)
		return
	}

	// Hvis limit er større enn antallet byer, bruk hele listen
	if limit > len(temp.Cities) {
		limit = len(temp.Cities)
	}

	// Legg til byene til country struct
	c.Cities = append(c.Cities, temp.Cities[:limit]...)

	// Sorter byene i alfabetisk rekkefølge
	sort.Strings(c.Cities)
}

// PrintCountry skriver ut landinformasjonen som JSON
func PrintCountry(w http.ResponseWriter, c Country) {
	// Struktur for å lagre landdataene som skal returneres
	var country struct {
		Name       string            `json:"name"`
		Continents []string          `json:"continents"`
		Languages  map[string]string `json:"languages"`
		Population int               `json:"population"`
		Borders    []string          `json:"borders"`
		Flag       string            `json:"flag"`
		Capital    []string          `json:"capital"`
		Cities     []string          `json:"cities"`
	}

	// Fyll ut strukturen med dataene fra Country-objektet
	country.Name = c.Name.Common
	country.Continents = c.Continents
	country.Languages = c.Languages
	country.Population = c.Population
	country.Borders = c.Borders
	country.Flag = c.Flag
	country.Capital = c.Capital
	country.Cities = c.Cities

	// Formaterer JSON med innrykk
	jsonCOUNTRY, err := json.MarshalIndent(country, "", "    ")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Skriv ut JSON-responsen
	fmt.Fprint(w, string(jsonCOUNTRY))
}
