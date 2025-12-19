FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o awesomeProject ./cmd/logitrack

FROM debian:bullseye-slim
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/awesomeProject ./awesomeProject
COPY config.yaml ./config.yaml
RUN chmod +x ./awesomeProject
EXPOSE 8080
CMD ["./awesomeProject"]