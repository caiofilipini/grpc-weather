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

var (
	EmptyResult = WeatherInfo{}
)

func elapsed(p WeatherProvider, start time.Time) {
	log.Printf("[%s] Request took %s\n", p.Name(), time.Since(start))
}
