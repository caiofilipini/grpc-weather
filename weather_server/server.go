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
	providers []*provider.WeatherProvider
}

func (s *server) registerProvider(p provider.WeatherProvider) {
	s.providers = append(s.providers, &p)
}

func (s server) queryProviders(q string) (*weather.WeatherResponse, error) {
	var responses []*provider.WeatherInfo
	var err error

	for _, p := range s.providers {
		prov := *p
		resp, e := prov.Query(q)
		if e != nil {
			err = e
		} else {
			log.Printf("[WeatherServer] Temperature obtained from %s is %.1f\n",
				prov.Name(), resp.Temperature)
			responses = append(responses, &resp)
		}
	}

	if len(responses) == 0 {
		return nil, err
	}

	return s.avg(responses), nil
}

func (s server) avg(responses []*provider.WeatherInfo) *weather.WeatherResponse {
	var sumTemp float64 = 0
	for _, r := range responses {
		if r.Found {
			sumTemp += r.Temperature
		}
	}
	avgTemp := sumTemp / float64(len(responses))

	return &weather.WeatherResponse{
		Temperature: avgTemp,
		Description: responses[0].Description,
		Found:       true,
	}
}

func (s server) CurrentConditions(ctx context.Context, req *weather.WeatherRequest) (*weather.WeatherResponse, error) {
	log.Println("[WeatherServer] Fetching weather information for", req.Location)
	defer elapsed(time.Now())

	return s.queryProviders(req.Location)
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

	weatherServer := &server{}
	weatherServer.registerProvider(provider.OpenWeatherMap{ApiKey: owmApiKey})
	weatherServer.registerProvider(provider.WeatherUnderground{ApiKey: wuApiKey})

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
