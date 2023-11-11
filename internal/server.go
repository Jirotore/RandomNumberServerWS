package internal

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"math/rand"
	"net/http"
	"net/netip"
	"sort"
	"sync"
	"time"
)

const (
	WebServerHostDefault   = "0.0.0.0"
	WebServerPortDefault   = "8080"
	WebServerAddrDefault   = WebServerHostDefault + ":" + WebServerPortDefault
	WebSocketServerDefault = "ws://" + WebServerAddrDefault + "/ws"
)

var (
	clientAddrs   = make([]netip.AddrPort, 0, 1000)
	uniqNumbers   = map[int64]struct{}{}
	clientAddrsMX = sync.Mutex{}
	uniqNumbersMX = sync.Mutex{}
)

func GeneratorInt64(conn *websocket.Conn) {
	logrus.Info("new connection")
	defer conn.Close()

	if err := checkDuplicateConn(conn.Request().RemoteAddr); err != nil {
		logrus.Errorf("cannot generate int64 number - %v\n", err)
		_ = conn.WriteClose(http.StatusForbidden)
		return
	}

	var request RequestModel[string]
	var response ResponseModel[int64]
	var err error

	for {
		err = websocket.JSON.Receive(conn, &request)

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			logrus.Errorf("cannot read data from client - %v\n", err)
			break
		}

		if len(request.Request) == 0 {
			logrus.Info("empty request from client")
			break
		}

		switch request.Request {
		case "get_number":
			response.Response = getNumber()
			err = websocket.JSON.Send(conn, &response)
			if err != nil {
				logrus.Errorf("cannot send data - %v\n", err)
				break
			}
		default:
			logrus.Errorf("unknown command - '%s'", request.Request)
		}
	}
}

func checkDuplicateConn(remoteAddr string) error {
	addrPort, err := netip.ParseAddrPort(remoteAddr)
	if err != nil {
		return errors.Join(errors.New("cannot parse address and port"), err)
	}

	filter := func(i int) bool {
		return clientAddrs[i].Addr().String() == addrPort.Addr().String() &&
			clientAddrs[i].Port() != addrPort.Port()
	}

	for !clientAddrsMX.TryLock() {
		time.Sleep(10 * time.Millisecond)
	}
	defer clientAddrsMX.Unlock()

	idx := sort.Search(len(clientAddrs), filter)

	if 0 <= idx && idx < len(clientAddrs) {
		return errors.New(fmt.Sprintf("brake new connection '%s' - already exists\n", addrPort.Addr().String()))
	}

	clientAddrs = append(clientAddrs, addrPort)
	sort.Slice(clientAddrs, func(i, j int) bool { return clientAddrs[i].Addr().Compare(clientAddrs[j].Addr()) <= 0 })

	return nil
}

func getNumber() int64 {
	for !uniqNumbersMX.TryLock() {
		time.Sleep(10 * time.Millisecond)
	}
	defer uniqNumbersMX.Unlock()

	rnd := rand.Int63()
	_, ok := uniqNumbers[rnd]
	for ok {
		rnd = rand.Int63()
		_, ok = uniqNumbers[rnd]
	}

	uniqNumbers[rnd] = struct{}{}

	return rnd
}
