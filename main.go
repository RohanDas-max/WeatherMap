package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func main() {

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("hello")
	})

	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		fmt.Println(city)
		data, err := query(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
	}
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	// ("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey)
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&apiid=" + apiConfig.OpenWeatherMapApiKey)
	if err != nil {
		return weatherData{}, err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	fmt.Println(d.Main, d.Name, d.Main.Kelvin)
	return d, nil

}
