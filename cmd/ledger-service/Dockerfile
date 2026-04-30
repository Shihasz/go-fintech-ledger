# Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build the specific ledger-service binary
RUN go build -o main cmd/ledger-service/main.go

# Run Stage
FROM alpine:3.19
WORKDIR /app

# Copy the compiled binary and the config file
COPY --from=builder /app/main .

EXPOSE 8081
CMD ["/app/main"]