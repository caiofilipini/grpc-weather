package main

import (
	"grcp/weather"
	"grcp/weather_server/provider"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	port = ":9000"
)

type server struct {
	provider provider.WeatherProvider
}

func (s server) mapResponse(i provider.WeatherInfo) *weather.WeatherResponse {
	return &weather.WeatherResponse{
		Temperature: i.Temperature,
		Description: i.Description,
		Found:       i.Found,
	}
}

func (s server) CurrentConditions(ctx context.Context, req *weather.WeatherRequest) (*weather.WeatherResponse, error) {
	log.Println("Fetching weather information for", req.Location)

	weatherInfo, err := s.provider.Query(req.Location)
	if err != nil {
		return nil, err
	}

	return s.mapResponse(weatherInfo), nil
}

func main() {
	conn, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}
	log.Println("Listening on", port)

	weatherServer := &server{
		provider: provider.OpenWeatherMap{
			ApiKey: os.Getenv("OPEN_WEATHER_MAP_API_KEY"),
		},
	}

	grpcServer := grpc.NewServer()
	weather.RegisterWeatherServer(grpcServer, weatherServer)
	grpcServer.Serve(conn)
}
