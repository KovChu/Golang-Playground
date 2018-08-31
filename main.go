package main

import (
	"strings"
	"encoding/json"
	"net/http"
)

func getAPIKey()(api_key string) {
	return "94582b9e9af7a138ef5e4c2101b395f9"
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", handleWeatherRequest)
	http.ListenAndServe(":8000", nil);
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!\n"))
}

func handleWeatherRequest(writer http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, error := query(city)

	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}else{
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(data)
	}
}

func query(city string) (weatherData, error) {
	resp, error := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + getAPIKey() + "&q=" + city);
	if error != nil {
		return weatherData{}, error
	}
	defer resp.Body.Close()

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, error
	}
	data.Main.Celius =  data.Main.Kelvin - 273
	return data, nil
}

type weatherData struct {
	Name string `json:"name"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	// Weather struct {
	// 	Weather string `json:"main"`
	// 	Detail string `json:"description"`
	// } `json:"weather"`
    Main struct {
		Kelvin float64 `json:"temp"`
		Celius float64
    } `json:"main"`
}