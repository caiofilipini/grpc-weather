# gRPC Weather

Simple gRPC example implemented in Go. It features a server and a client for a service that provides information about the weather, and it relies on [openweathermap.org's API](http://openweathermap.org/api) for that.

## Building and running

You will need to [install protoc](https://github.com/google/protobuf/blob/master/INSTALL.txt) and the [protoc-gen-go](https://github.com/golang/protobuf) plugin in order to generate server and client stubs.

Additionally, you will need an [openweathermap.org API key](http://openweathermap.org/appid).

Once that's all in place, you can build the server:

```sh
$ make build-server
```

And then run it:

```sh
$ OPEN_WEATHER_MAP_API_KEY="s3cr3+" ./weather_server/server
```

If all goes well, you should see a message like this (assuming the server is running in the default port (i.e. `9000`):

```sh
2015/11/18 17:32:59 Listening on :9000
```

Similarly, you can build the client:

```sh
$ make build-client
```

And then run it, providing a location:

```sh
./weather_client/client Berlin, Germany
```

An example output:

```sh
2015/11/18 17:33:03 It's currently 9.5°C, broken clouds in Berlin, Germany
```
