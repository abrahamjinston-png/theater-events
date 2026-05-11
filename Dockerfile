FROM golang:1.26-bookworm

RUN apt-get update && apt-get install -y \
    chromium \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -tags netgo -ldflags="-s -w" -o app ./cmd/scraper

ENV CHROME_PATH=/usr/bin/chromium

CMD ["./app"]