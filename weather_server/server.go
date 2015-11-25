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
	"github.com/caiofilipini/grpc-weather/weather_server/providers"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	defaultPort = 9000
)

type server struct {
	providers *providers.WeatherProviders
}

func (s server) CurrentConditions(ctx context.Context, req *weather.WeatherRequest) (*weather.WeatherResponse, error) {
	log.Println("[WeatherServer] Fetching weather information for", req.Location)
	defer elapsed(time.Now())

	response, err := s.providers.Query(req.Location)
	if err != nil {
		return nil, err
	}

	return &weather.WeatherResponse{
		Temperature: response.Temperature,
		Description: response.Description,
		Found:       response.Found,
	}, nil
}

func main() {
	owmApiKey := strings.TrimSpace(os.Getenv("OPEN_WEATHER_MAP_API_KEY"))
	wuApiKey := strings.TrimSpace(os.Getenv("WEATHER_UNDERGROUND_API_KEY"))
	if owmApiKey == "" {
		log.Fatal("Missing API key for OpenWeatherMap")
	}
	if wuApiKey == "" {
		log.Fatal("Missing API key for Weather Underground")
	}

	weatherServer := &server{providers: &providers.WeatherProviders{}}
	weatherServer.providers.Register(providers.OpenWeatherMap{ApiKey: owmApiKey})
	weatherServer.providers.Register(providers.WeatherUnderground{ApiKey: wuApiKey})

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
		log.Fatalf("Failed to listen on %d: %v", port, err)
	}

	log.Println("[WeatherServer] Listening on", port)
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
	log.Printf("[WeatherServer] Request took %s\n", time.Since(start))
}
