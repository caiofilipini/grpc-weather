package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
)

const (
	openWeatherMapUrl = "http://api.openweathermap.org/data/2.5/weather"
)

type WeatherProvider interface {
	Query(string) (WeatherInfo, error)
}

type OpenWeatherMap struct {
	ApiKey string
}

type WeatherInfo struct {
	Temperature float64
	Description string
	Found       bool
}

func (p OpenWeatherMap) Query(q string) (WeatherInfo, error) {
	body, err := p.get(q)
	if err != nil {
		return WeatherInfo{}, err
	}

	var result openWeatherMapResult
	json.Unmarshal(body, &result)

	if result.found() {
		return WeatherInfo{
			Temperature: result.toCelcius(),
			Description: result.description(),
			Found:       true,
		}, nil
	} else {
		return WeatherInfo{Found: false}, nil
	}
}

func (p OpenWeatherMap) get(q string) ([]byte, error) {
	queryUrl := fmt.Sprintf("%s?q=%s&appid=%s", openWeatherMapUrl, url.QueryEscape(q), p.ApiKey)

	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected response: %s", resp.Status)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
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

func (r openWeatherMapResult) toCelcius() float64 {
	return math.Floor(r.Main.Kelvin-273.15) + 0.5
}

func (r openWeatherMapResult) description() string {
	return r.Weather[0].Description
}

func (r openWeatherMapResult) found() bool {
	return len(r.Weather) > 0
}
