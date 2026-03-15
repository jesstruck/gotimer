FROM golang:1.22-bookworm AS builder

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/time-tracker-backend ./main.go

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -r app \
    && useradd -r -g app app \
    && mkdir -p /data \
    && chown -R app:app /data

COPY --from=builder /out/time-tracker-backend /usr/local/bin/time-tracker-backend

ENV TIME_ENTRIES_DB_PATH=/data/time_entries.db

USER app
EXPOSE 8080

CMD ["time-tracker-backend"]
