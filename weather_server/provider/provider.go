package provider

type WeatherProvider interface {
	Query(string) (WeatherInfo, error)
}

type WeatherInfo struct {
	Temperature float64
	Description string
	Found       bool
}
