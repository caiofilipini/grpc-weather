package providers

import (
	"log"
	"time"
)

type WeatherProvider interface {
	Query(string) (WeatherInfo, error)
	Name() string
}

type WeatherInfo struct {
	Temperature float64
	Description string
	Found       bool
}

type WeatherProviders struct {
	providers []*WeatherProvider
}

func (wps *WeatherProviders) Register(provider WeatherProvider) {
	wps.providers = append(wps.providers, &provider)
}

func (wps *WeatherProviders) Query(q string) (*WeatherInfo, error) {
	var responses []*WeatherInfo
	var err error

	for _, p := range wps.providers {
		prov := *p
		resp, e := prov.Query(q)
		if e != nil {
			err = e
		} else {
			log.Printf("[WeatherProviders] Temperature obtained from %s is %.1f\n",
				prov.Name(), resp.Temperature)
			responses = append(responses, &resp)
		}
	}

	if len(responses) == 0 {
		return nil, err
	}

	return wps.avg(responses), nil
}

func (wps *WeatherProviders) avg(responses []*WeatherInfo) *WeatherInfo {
	var sumTemp float64 = 0
	for _, r := range responses {
		if r.Found {
			sumTemp += r.Temperature
		}
	}
	avgTemp := sumTemp / float64(len(responses))

	return &WeatherInfo{
		Temperature: avgTemp,
		Description: responses[0].Description,
		Found:       true,
	}
}

var (
	EmptyResult = WeatherInfo{}
)

func elapsed(p WeatherProvider, start time.Time) {
	log.Printf("[%s] Request took %s\n", p.Name(), time.Since(start))
}
