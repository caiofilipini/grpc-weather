package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/caiofilipini/grpc-weather/weather"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	port *int
	host *string
)

func init() {
	host = flag.String("s", "localhost", "server host")
	port = flag.Int("p", 9000, "server port")

	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		log.Fatalf("Missing location parameter")
	}

	serverAddr := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to %s: %v", serverAddr, err)
	}
	defer conn.Close()

	location := strings.Join(flag.Args(), " ")
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
