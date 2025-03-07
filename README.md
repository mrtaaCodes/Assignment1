# COUNTRY INFORMATION API
## Assignment 1

Documantation for assignment 1. This API allows you to retrieve essential data about countries, such as general information, population statistics over specific periods, and real-time API status checks.

---



## API Endpoints

### 1. **Check API Status**  
**GET** `/status`  
This endpoint provides information about the current status of the API.

URL:
https://assignment1-cryh.onrender.com/countryinfo/v1/status
#### Example:


```bash
/status
{
    "CountriesNowStatus": "200 OK",
    "CountriesRestStatus": "200 OK",
    "Version": "v1",
    "Uptime": 31
}
```
This response indicates that all systems are functioning properly, with a 200 OK status confirming that the service is operational.



# 2. Retrieve Country Information
GET /info/{ISO2-country_code}?limit=integer

This endpoint allows you to retrieve detailed information about a specific country, such as its name, capital, population, languages spoken, borders, and a list of cities.
The limit parameter (optional) can be used to specify the maximum number of cities to include in the response.


URL:
https://assignment1-cryh.onrender.com/countryinfo/v1/info/gb
Example:
```
/info/no?limit=5

{
    "name": "United Kingdom",
    "continents": [
        "Europe"
    ],
    "languages": {
        "eng": "English"
    },
    "population": 67215293,
    "borders": [
        "IRL"
    ],
    "flag": "🇬🇧",
    "capital": [
        "London"
    ],
    "cities": [
        "Abberton",
        "Abbots Langley",
        "Aberaeron",
        "Aberchirder",
        "Abercynon",
        "Aberdare",
        "Aberdeen",
        "Aberfeldy",
        "Aberford",
        "Aberfoyle"
    ]
}
```



# 3. Retrieve Population Data
GET /population/{ISO2-country_code}?limit="startYear-endYear"

This endpoint provides population statistics for a country over a specified range of years. The limit parameter allows you to filter the data by a specific time range (e.g., "2000-2005").

URL:
https://assignment1-cryh.onrender.com/countryinfo/v1/population/gb?limit=2000-2005
Example:
```
/population/gb?limit=2000-2005

{
    "mean": 59569892,
    "values": [
        {
            "year": 2000,
            "value": 58892514
        },
        {
            "year": 2001,
            "value": 59119673
        },
        {
            "year": 2002,
            "value": 59370479
        },
        {
            "year": 2003,
            "value": 59647577
        },
        {
            "year": 2004,
            "value": 59987905
        },
        {
            "year": 2005,
            "value": 60401206
        }
    ]
}
```




#Summary of Endpoints
Check API Status

Endpoint: /status
Description: Retrieves the current status of the API.
Example: 
Indicates whether the API is functioning properly, showing a 200 OK status.
Retrieve Country Information

Endpoint: /info/{ISO2-country_code}?limit=integer
Description: Provides detailed information about a country (e.g., name, capital, population, languages, borders, cities).
Example:
Returns country data, including a list of cities and other relevant information.
Retrieve Population Data

Endpoint: /population/{ISO2-country_code}?limit="startYear-endYear"
Description: Provides population statistics for a country over a specified range of years.
Example: 
Returns population data for the specified years, including the mean population and individual yearly values.

