# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Ensure go.mod is up to date and build static binary
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o cli ./cmd/cli

# Runtime stage - minimal image (~2MB)
FROM scratch

COPY --from=builder /build/cli /cli

ENTRYPOINT ["/cli"]
