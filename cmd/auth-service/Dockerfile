# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
# Note: We copy the whole internal folder so the build can access shared packages
COPY cmd/auth-service ./cmd/auth-service
COPY internal ./internal

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/auth-service

# Final Stage (Distroless)
# We use gcr.io/distroless/static for minimal footprint and security
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/auth-service .

# Run as non-root user (security best practice)
USER 65532:65532

ENTRYPOINT ["/auth-service"]
