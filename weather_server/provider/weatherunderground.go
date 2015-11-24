package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	wuUrlTemplate = "http://api.wunderground.com/api/%s/conditions"

	errKeyNotFound   = "keynotfound"
	errQueryNotFound = "querynotfound"
)

type WeatherUnderground struct {
	ApiKey string
}

func (p WeatherUnderground) Query(q string) (WeatherInfo, error) {
	result, err := p.getAsJSON(p.urlFor(q))
	if err != nil {
		return EmptyResult, err
	}

	if apiErr := result.Response.Error.Type; apiErr != "" {
		if apiErr == errQueryNotFound {
			return EmptyResult, nil
		}
		return EmptyResult, fmt.Errorf("Error querying Weather Underground: %s", apiErr)
	}

	if len(result.Response.Results) > 0 {
		var location string
		for _, r := range result.Response.Results {
			if strings.Contains(q, r.City) {
				location = r.Location
				break
			}
		}

		if location != "" {
			result, err = p.getAsJSON(p.urlFor(location))
			if err != nil {
				return EmptyResult, nil
			}
		}
	}

	return result.asWeatherInfo(), nil
}

func (p WeatherUnderground) urlFor(q string) string {
	baseUrl := fmt.Sprintf(wuUrlTemplate, p.ApiKey)

	var queryPart string
	if strings.Contains(q, "/q/") {
		queryPart = q
	} else {
		queryPart = fmt.Sprintf("/q/%s", strings.Replace(url.QueryEscape(q), "+", "%20", 1))
	}

	return fmt.Sprintf("%s%s.json", baseUrl, queryPart)
}

func (p WeatherUnderground) getAsJSON(q string) (weatherUndergroundResult, error) {
	var result weatherUndergroundResult

	body, err := p.get(q)
	if err != nil {
		return result, err
	}

	json.Unmarshal(body, &result)

	return result, nil
}

func (p WeatherUnderground) get(queryUrl string) ([]byte, error) {
	log.Println("Querying WU with", queryUrl)

	resp, err := http.Get(queryUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		respErr := fmt.Errorf("Unexpected response: %s", resp.Status)
		log.Println("Request failed:", respErr)
		return nil, respErr
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type weatherUndergroundResult struct {
	Response struct {
		Results []struct {
			City     string `json:"city"`
			Location string `json:"l"`
		} `json:"results,omitempty"`
		Error struct {
			Type string `json:"type"`
		} `json:"error,omitempty"`
	} `json:"response"`
	CurrentObservation struct {
		TempC   float64 `json:"temp_c"`
		Weather string  `json:"weather"`
	} `json:"current_observation,omitempty"`
}

func (r weatherUndergroundResult) asWeatherInfo() WeatherInfo {
	return WeatherInfo{
		Temperature: r.CurrentObservation.TempC,
		Description: r.CurrentObservation.Weather,
		Found:       true,
	}
}
