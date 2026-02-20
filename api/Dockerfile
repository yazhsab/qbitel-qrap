FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app
COPY go.work go.work.sum* ./
COPY api/go.mod api/go.sum ./api/
COPY shared/go/go.mod shared/go/go.sum ./shared/go/
RUN cd api && go mod download

COPY api/ ./api/
COPY shared/go/ ./shared/go/
RUN cd api && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /qrap-api ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /qrap-api /usr/local/bin/qrap-api

EXPOSE 8083
ENTRYPOINT ["qrap-api"]
