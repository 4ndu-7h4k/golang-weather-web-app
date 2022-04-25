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

var (
	tmpl *template.Template
 	setEnv = godotenv.Load(".env");
)
type weatherData struct {
	Name      string 	`json:"name"`
	StatuCode uint16  	`json:"cod"`
	TimeZone  int32 	`json:"timezone"`

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

type portData struct{
	Port string
	Domain string
}


func init(){

	if(setEnv != nil){ 
		log.Fatal("[ERROR] : failed to load `.env` file")
	}

	_,err := os.LookupEnv("PORT")
	if(!err){

		os.Setenv("PORT","80")
		log.Println("[WARNING] : env `PORT` not set , Default running port on 80 ")
	}

	
	_,err = os.LookupEnv("DOMAIN")
	if(!err){
		os.Setenv("DOMAIN","localhost")
		log.Println("[WARNING] : env `DOMAIN` not set , Deafult set to localhost ")
	}

	_,err = os.LookupEnv("OWMAPI")
	if(!err){
		log.Fatal("[ERROR] : env `OWMAPI` not set , Check enviornment variable list")
	}

	log.Println("`.env` loaded successfully")

}

func homePage(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	domain := os.Getenv("DOMAIN")
	data := portData{
		Port : port,
		Domain : domain,
	}
	tmpl.Execute(w, data)
}

func query(city string) (weatherData, error) {
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
}

//to get the time by offset
func getTimeOffset(offset int64) string {

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).Add(time.Second * time.Duration(offset))
	return now.Format("03:04 PM") 

}

//function to return the time
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(currTimeJson)

}
func getCity(w http.ResponseWriter, r *http.Request) {

	d := &cityData{
		Cities: []string{
			"Kochi,IN",
			"Kollam,IN",
			"Kozhikode,IN",
			"Kasargode,IN",
			"Chennai,IN",
			"Mumbai,IN",
			"London,UK",
			"Pathanamthitta,IN",
		},
	}
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
	port := os.Getenv("PORT")
	log.Println("Staring http server at port : "+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
