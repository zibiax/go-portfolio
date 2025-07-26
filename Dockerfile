FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go run download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["/docker/gs-ping"]
