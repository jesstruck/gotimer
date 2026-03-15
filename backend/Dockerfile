FROM golang:1.22-alpine AS builder

WORKDIR /src

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/time-tracker-backend ./main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /data \
    && chown -R app:app /data

COPY --from=builder /out/time-tracker-backend /usr/local/bin/time-tracker-backend

ENV TIME_ENTRIES_DB_PATH=/data/time_entries.db

USER app
EXPOSE 8080

CMD ["time-tracker-backend"]
