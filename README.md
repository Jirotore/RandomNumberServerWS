# RandomNumberServerWS
websocket server receive random int64 number for websocket client

## Build Applications

### [Docker]
run command `docker-compose build` and build both apps: websocket server and client

### [Linux Bash]
build websocket server - `go build -o server main.go`

build websocket client - `go build -o clientWS client/client.go`

## Run apps
### [Linux Bash]
* websocket server - `/bin/bash ./server`
* websocket client - `/bin/bash ./clientWS`

## Environments

server envs:
* `APP_WS_SERVER_HOST` - ip address of server (default: `0.0.0.0`)
* `APP_WS_SERVER_PORT` - port number of server (default: `8080`)

client envs:
* `APP_WS_SERVER_URL` - url websocket server (default: `ws://0.0.0.0:8080/ws`)
