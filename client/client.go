package main

import (
	"RandomWS/internal"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"os"
)

func main() {
	wsServerUrl, ok := os.LookupEnv("APP_WS_SERVER_URL")
	if !ok {
		wsServerUrl = internal.WebSocketServerDefault
	}

	conn, err := websocket.Dial(wsServerUrl, "", wsServerUrl)
	if err != nil {
		logrus.Errorf("cannot connect to server - %v\n", err)
		return
	}
	defer conn.Close()

	var request internal.RequestModel[string]
	var response internal.ResponseModel[int64]

	for response.ErrorMsg == "" && err == nil {
		request.Request = "get_number"
		err = websocket.JSON.Send(conn, &request)
		if err != nil {
			logrus.Errorf("write error - %v\n", err)
		}

		err = websocket.JSON.Receive(conn, &response)
		if err != nil {
			logrus.Errorf("cannot read data - %v\n", err)
			break
		}

		if response.ErrorMsg != "" {
			logrus.Errorf("got error response - %s", response.ErrorMsg)
			break
		}

		logrus.Info("random number = ", response.Response)
	}
}
