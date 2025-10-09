FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy custodian-data-adapter-go dependency
COPY custodian-data-adapter-go/ ./custodian-data-adapter-go/

# Copy custodian-simulator-go files
COPY custodian-simulator-go/go.mod custodian-simulator-go/go.sum ./custodian-simulator-go/
WORKDIR /build/custodian-simulator-go
RUN go mod download

# Copy source and build
COPY custodian-simulator-go/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o custodian-simulator ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates wget
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app
COPY --from=builder /build/custodian-simulator-go/custodian-simulator /app/custodian-simulator
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080 50051

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8084/api/v1/health || exit 1

CMD ["./custodian-simulator"]