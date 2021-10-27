package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Weather struct {
	Base   string `json:"base"`
	Clouds struct {
		All int64 `json:"all"`
	} `json:"clouds"`
	Cod   int64 `json:"cod"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Dt   int64 `json:"dt"`
	Id   int64 `json:"id"`
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Humidity  int64   `json:"humidity"`
		Pressure  int64   `json:"pressure"`
		Temp      float64 `json:"temp"`
		TempMax   float64 `json:"temp_max"`
		TempMin   float64 `json:"temp_min"`
	} `json:"main"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		Id      int64  `json:"id"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
		Type    int64  `json:"type"`
	} `json:"sys"`
	Timezone   int64 `json:"timezone"`
	Visibility int64 `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Id          int64  `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Deg   float64  `json:"deg"`
		Gust  float64 `json:"gust"`
		Speed float64  `json:"speed"`
	} `json:"wind"`
}


func (w *Weather) GetWeather(c *gin.Context) {
	w.Init()
	c.IndentedJSON(http.StatusOK, w)
}

func (w *Weather) ReadFile() []byte {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?id=4671240&units=imperial&appid=%v", viper.GetString("APIKEY")))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func (w *Weather) Init() {
	err := json.Unmarshal(w.ReadFile(), w)
	if err != nil {
		log.Fatal(err)
	}
}


func main() {
	var weather Weather
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders: []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	router.GET("/weather", weather.GetWeather)
	if err := router.Run(":"+os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}
