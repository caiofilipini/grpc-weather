gen-proto:
	protoc --go_out=plugins=grpc:. weather/weather.proto

build-server: gen-proto
	go build -o weather_server/server weather_server/server.go

install-server: gen-proto
	go install github.com/caiofilipini/grpc-weather/weather_server

build-client: gen-proto
	go build -o weather_client/client weather_client/client.go

install-client: gen-proto
	go install github.com/caiofilipini/grpc-weather/weather_client

clean:
	rm -f weather/*.pb.go
	rm -f weather_server/server
	rm -f weather_client/client
