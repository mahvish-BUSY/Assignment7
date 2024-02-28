package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type WeatherData struct {
	Name    string   `json:"name"`
	Weather string   `json:"weather"`
	Status  []string `json:"status"`
}
type WeatherApiResponse struct {
	Page       int           `json:"page"`
	PerPage    int           `json:"per_page"`
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
	Data       []WeatherData `json:"data"`
}

func retrieveNumber(str string) (float64, error) {
	re := regexp.MustCompile(`\d*\.?\d+`)
	match := re.FindString(str)

	//convert it into number
	num, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, err
	}
	return num, nil

}
func fetchAllData(apiURL string) ([]interface{}, error) {
	//var to hold all data
	var alldata []interface{}
	page := 1

	for {
		//making a get request to api end point
		url := apiURL + "&page=" + strconv.Itoa(page)
		response, err := http.Get(url)
		if err != nil {
			response.Body.Close()
			return nil, err
		}

		//parsing the data into struct
		var weatherResponse WeatherApiResponse
		err = json.NewDecoder(response.Body).Decode(&weatherResponse)
		if err != nil {
			return nil, err
		}
		response.Body.Close()

		for _, weatherData := range weatherResponse.Data {
			name := weatherData.Name
			degree, err := retrieveNumber(weatherData.Weather)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			wind, err := retrieveNumber(weatherData.Status[0])
			if err != nil {
				fmt.Println("Error: ", err)
			}
			humidity, err := retrieveNumber(weatherData.Status[1])
			if err != nil {
				fmt.Println("Error: ", err)
			}
			alldata = append(alldata, []interface{}{name, degree, wind, humidity})

		}

		//check for next page, if available proceed else break
		//break statement
		if page > weatherResponse.TotalPages {
			break
		}
		page++
	}
	//return response
	return alldata, nil
}
func main() {

	//api end point
	apiURL := "https://jsonmock.hackerrank.com/api/weather/search?name="
	fmt.Println("Enter the name string/substring to be searched")
	var name string
	fmt.Scanln(&name)
	apiURL = apiURL + name
	
	allData, err := fetchAllData(apiURL)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Total weather data fetched : %d\n\n", len(allData))

	for _, data := range allData {
		if cityData, ok := data.([]interface{}); ok {
			for _, value := range cityData {
				fmt.Printf("\t%v", value)
			}
		}
		fmt.Println()
	}

}
