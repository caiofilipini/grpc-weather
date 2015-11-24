package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caiofilipini/grpc-weather/weather"
	"github.com/caiofilipini/grpc-weather/weather_server/provider"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	defaultPort = 9000
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
	defer elapsed(time.Now())

	weatherInfo, err := s.provider.Query(req.Location)
	if err != nil {
		return nil, err
	}

	return s.mapResponse(weatherInfo), nil
}

func main() {
	owmApiKey := strings.TrimSpace(os.Getenv("OPEN_WEATHER_MAP_API_KEY"))
	if owmApiKey == "" {
		log.Fatal("Missing API key for OpenWeatherMap")
	}

	weatherServer := &server{
		provider: provider.OpenWeatherMap{
			ApiKey: owmApiKey,
		},
	}

	conn := listen()
	grpcServer := grpc.NewServer()
	weather.RegisterWeatherServer(grpcServer, weatherServer)
	grpcServer.Serve(conn)
}

func listen() net.Listener {
	port := assignPort()
	listenAddr := fmt.Sprintf(":%d", port)

	conn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	log.Println("Listening on", port)
	return conn
}

func assignPort() int {
	if p := os.Getenv("PORT"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Invalid port %s", p)
		}
		return port
	}
	return defaultPort
}

func elapsed(start time.Time) {
	log.Printf("Request took %s\n", time.Since(start))
}
