# The JS bundle and the Go binary are built on the runner's own platform and
# cross-compiled; only the final stage is per-architecture (amd64 + arm64).
FROM --platform=$BUILDPLATFORM node:22-slim AS frontend

WORKDIR /app/svelte-components

COPY svelte-components/package.json svelte-components/package-lock.json ./
RUN npm ci

COPY svelte-components/ ./
COPY static/css ../static/css
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.25 AS build

ARG TARGETOS TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /portfolio .

FROM debian:bookworm-slim

# Copied rather than apt-installed, so this stage has no RUN step and needs
# no emulation when built for another architecture.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

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
