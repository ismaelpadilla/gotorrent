FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY clients ./clients
COPY cmd ./cmd
COPY interfaces ./interfaces
COPY ui ./ui
RUN go build -o /gotorrent

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /gotorrent /gotorrent

USER nonroot:nonroot

ENTRYPOINT ["/gotorrent"]
