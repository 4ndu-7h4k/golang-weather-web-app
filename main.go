package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/joho/godotenv"
)

var tmpl *template.Template

type weatherData struct {
	Name      string `json:"name"`
	StatuCode uint8  `json:"cod"`
	TimeZone  int32 `json:"timezone"`
	Cloud     []struct {
		Icon        string `json:"icon"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Celsius  float64 `json:"temp"`
		Humidity float32 `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float32 `json:"speed"`
	} `json:"wind"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
}
type timeData struct {
	Time string `json:"time"`
}
type cityData struct {
	Cities []string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("hello from go!\n"))
	tmpl.Execute(w, nil)
}
func query(city string) (weatherData, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	apiKey := os.Getenv("OWMAPI")
	resp, err := http.Get("http://api.openweatherMap.org/data/2.5/weather?APPID=" + apiKey + "&q=" + city + "&units=metric")
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	log.Println(d)
	return d, nil
}
func getWeather(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
	log.Println(data)
}

//to get the time by offset
func getTimeOffset(offset int64) string {

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).Add(time.Second * time.Duration(offset))
	return now.Format("03:04 PM") // UTC

}

//func
func getTime(w http.ResponseWriter, r *http.Request) {

	offsetString := strings.SplitN(r.URL.Path, "/", 3)[2]
	offset, err := strconv.ParseInt(offsetString, 0, 32)
	if err != nil {
		log.Println("Unable cast offset '" + offsetString + "'")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currTime := getTimeOffset(offset)
	currTimeJson := &timeData{
		Time: currTime,
	}
	data, err := json.Marshal(currTimeJson)
	if err != nil {
		log.Println("Unable to get json")
	}
	log.Println(data)
	log.Println("Time API triggered : Response " + currTimeJson.Time)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(currTimeJson)

}
func getCity(w http.ResponseWriter, r *http.Request) {

	d := &cityData{
		Cities: []string{
			"Kochi,IN",
			"Kollam,IN",
			"Kozhikode,IN",
			"Konni,IN",
			"Kasargode,IN",
			"Chennai,IN",
			"Mumbai,IN",
			"London,UK",
			"Pathanamthitta,IN",
		},
	}
	log.Println(d)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(d)
}
func main() {
	fs := http.FileServer(http.Dir("./static"))
	tmpl = template.Must((template.ParseFiles("template/index.html")))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/time/", getTime)
	http.HandleFunc("/getCity/", getCity)
	http.HandleFunc("/weather/", getWeather)
	log.Println("Staring http server at port : 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
