FROM golang:1.26.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/peakcloud-api \
    ./cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates curl \
    && addgroup -S peakcloud \
    && adduser -S -G peakcloud peakcloud

WORKDIR /app

COPY --from=builder /out/peakcloud-api /app/peakcloud-api

USER peakcloud

EXPOSE 8080

ENTRYPOINT ["/app/peakcloud-api"]
