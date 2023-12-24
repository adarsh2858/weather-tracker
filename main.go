package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
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

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	if err := json.Unmarshal(bytes, &c); err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}

	resp, err := http.Get("http://openweathermap.com/temp?APPID=" + apiConfig.OpenWeatherMapAppID + "&q=" + city)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
