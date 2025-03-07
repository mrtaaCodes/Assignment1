# COUNTRY INFORMATION API
# Assignment 1

Welcome to the **Oblig_1** API documentation. This service provides a comprehensive interface to retrieve essential data about countries, including general country information, population statistics over specified time periods, and real-time API status checks.

---

## API Endpoints

### 1. Retrieve Country Information  
**GET** `/info/{ISO2-country_code}?limit=integer`  
This endpoint allows you to obtain detailed information about a specific country, including its name, capital, population, languages spoken, borders, and a list of cities.

- **Parameters:**
  - `limit` (optional): Specifies the maximum number of cities to include in the response.

#### Example Request:
```bash
/info/no?limit=5
```
https://assignment1-cryh.onrender.com/countryinfo/v1/status
```bash
{
    "CountriesNowStatus": "200 OK",
    "CountriesRestStatus": "200 OK",
    "Version": "v1",
    "Uptime": 31
}
```

