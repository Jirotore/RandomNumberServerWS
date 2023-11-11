FROM golang:1.20-alpine as build
WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY internal internal
COPY client client
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o clientWS client/client.go

FROM alpine:3.18.4 as server
WORKDIR /app
COPY --from=build /build/server /app/server
ENTRYPOINT /app/server

FROM alpine:3.18.4 as client
WORKDIR /app
COPY --from=build /build/clientWS /app/client
ENTRYPOINT /app/client
