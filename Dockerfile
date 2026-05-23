FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/impulse ./cmd/impulse

FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /out/impulse /app/impulse
COPY config.json /app/config.json
COPY events /app/events

RUN chown -R app:app /app

USER app

CMD ["./impulse", "config.json", "events"]
