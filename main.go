package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

var ()

func main() {
	// printPath()
	// practiceArray()
	// http.HandleFunc("/", hello)
	// timeTest()
	// stringlen()
	// http.HandleFunc("/weather/", handleWeatherRequest)
	mw := multiWeatherProvider{openWeatherMap{apiKey: "94582b9e9af7a138ef5e4c2101b395f9"}}
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		temp, err := mw.temperature(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string] interface{}{
			"city": city,
			"temp": temp,
			"took": time.Since(begin).String(),
		})
	})
	http.ListenAndServe(":8000", nil)
}

func timeTest() {
}

func printPath() {
	fmt.Println("path of executable " + os.Args[0])
}

func practiceArray() {
	scores := make([]int, 10)

	for i := 0; i < 10; i++ {
		scores[i] = int(rand.Int31n(1000))
	}
	fmt.Println(scores)
	sort.Ints(scores)
	worst := make([]int, 5)
	copy(worst, scores[:3])
	fmt.Println(scores)
	fmt.Println(worst)
}

func stringlen() {
	fmt.Println(len("hi"))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!\n"))
}

func handleWeatherRequest(writer http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	m := openWeatherMap {apiKey: "94582b9e9af7a138ef5e4c2101b395f9"}
	data, error := m.query(city)

	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	} else {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(data)
	}
}

func (w openWeatherMap) query(city string) (weatherData, error) {
	fmt.Println("Getting weather info for " + city)
	resp, error := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if error != nil {
		return weatherData{}, error
	}
	defer resp.Body.Close()

	var data weatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherData{}, error
	}
	data.temperature.Celius = data.temperature.Kelvin - 273.15
	return data, nil
}

type weatherProvider interface {
	temperature(city string) (float64, error)
}

type openWeatherMap struct {
	apiKey string
}

func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + w.apiKey + "&q=" + city)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}

	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

func (w multiWeatherProvider) temperature(city string) (float64, error) {
	sum := 0.0

	for _, provider := range w {
		k, err := provider.temperature(city)
		if err != nil {
			return 0, err
		}

		sum += k
	}

	return sum / float64(len(w)), nil
}

type multiWeatherProvider []weatherProvider

type weatherData struct {
	Name string `json:"name"`
	// Weather struct {
	// 	Weather string `json:"main"`
	// 	Detail string `json:"description"`
	// } `json:"weather"`
	coord       `json:"coord"`
	temperature `json:"main"`
}

type coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type temperature struct {
	Kelvin float64 `json:"temp"`
	Celius float64
}
