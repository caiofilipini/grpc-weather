package main

import (
	"grcp/weather"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	serverAddr = "localhost:9000"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("Missing location parameter")
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", serverAddr, err)
	}
	defer conn.Close()

	location := strings.Join(os.Args[1:], " ")
	client := weather.NewWeatherClient(conn)
	req := weather.WeatherRequest{Location: location}

	resp, err := client.CurrentConditions(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	if resp.Found {
		log.Printf("It's currently %.1f°C, %s in %s\n", resp.Temperature, resp.Description, req.Location)
	} else {
		log.Printf("Could not find weather information for %s\n", req.Location)
	}
}
