FROM node:22-slim AS frontend

WORKDIR /app/svelte-components

COPY svelte-components/package.json svelte-components/package-lock.json ./
RUN npm ci

COPY svelte-components/ ./
COPY static/css ../static/css
RUN npm run build

FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /portfolio .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /portfolio ./portfolio
COPY templates ./templates
COPY static ./static
COPY --from=frontend /app/static/js/bundle.js /app/static/js/bundle.css ./static/js/
COPY robots.txt ./robots.txt

ENV BLOG_DB_PATH=/app/data/blog.db
VOLUME ["/app/data"]

EXPOSE 5000

CMD ["./portfolio"]
