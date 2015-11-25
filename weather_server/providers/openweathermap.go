package providers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"
)

const (
	openWeatherMapUrl  = "http://api.openweathermap.org/data/2.5/weather"
	openWeatherMapName = "OpenWeatherMap"
)

type OpenWeatherMap struct {
	ApiKey string
}

func (p OpenWeatherMap) Name() string {
	return openWeatherMapName
}

func (p OpenWeatherMap) Query(q string) (WeatherInfo, error) {
	defer elapsed(p, time.Now())

	queryUrl := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapUrl, url.QueryEscape(q), p.ApiKey)
	body, err := httpClient.get(queryUrl)
	if err != nil {
		return EmptyResult, err
	}

	var result openWeatherMapResult
	json.Unmarshal(body, &result)

	return result.asWeatherInfo(), nil
}

type openWeatherMapResult struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
	Sys struct {
		Country string `json:"country"`
	} `json:"sys"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func (r openWeatherMapResult) asWeatherInfo() WeatherInfo {
	if r.found() {
		return WeatherInfo{
			Temperature: r.toCelcius(),
			Description: r.description(),
			Found:       true,
		}
	}
	return EmptyResult
}

func (r openWeatherMapResult) toCelcius() float64 {
	return math.Floor(r.Main.Kelvin-273.15) + 0.5
}

func (r openWeatherMapResult) description() string {
	return r.Weather[0].Description
}

func (r openWeatherMapResult) found() bool {
	return len(r.Weather) > 0
}
