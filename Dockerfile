FROM golang:1.23.2
FROM node:22 AS build

WORKDIR /app

COPY package.json ./
COPY package.lock ./

COPY go.mod go.sum ./

RUN go run download && npm install

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["/docker/gs-ping"]
