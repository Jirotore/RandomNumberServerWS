package main

import (
	"RandomWS/internal"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
)

func main() {
	host, ok := os.LookupEnv("APP_WS_SERVER_HOST")
	if !ok {
		host = internal.WebServerHostDefault
	}

	port, ok := os.LookupEnv("APP_WS_SERVER_PORT")
	if !ok {
		port = internal.WebServerPortDefault
	}

	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(internal.GeneratorInt64))
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux))
}
