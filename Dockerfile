# ---------- Build Stage ----------
FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

# Copy go.mod & go.sum first for layer caching
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the rest of the code
COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api main.go

# ---------- Final Stage ----------
FROM alpine:3.18
WORKDIR /app

# Copy binary from builder
COPY --from=builder /api /app/api

# Add runtime dependencies
RUN apk add --no-cache git curl ca-certificates tzdata

# Set env and user
ENV APP_ENV=production
EXPOSE 8080
USER nobody:nobody

CMD ["/app/api"]