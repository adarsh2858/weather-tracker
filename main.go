package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type apiConfig struct {
	OpenWeatherMapAppID string `json:"open_weather_map_app_id"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		TempInKelvin string `json:"temp"`
	} `json:"main"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from go!\n"))
}

func loadApiConfig() apiConfig {
	return apiConfig{}
}

func query(city string) (weatherData, error) {
	apiConfig := loadApiConfig()

	resp, err := http.Get("http://openweathermap.com/temp?APPID=" + apiConfig.OpenWeatherMapAppID + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	json.NewDecoder(resp.Body).Decode(&d)

	return weatherData{}, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		_, _ = query(city)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
