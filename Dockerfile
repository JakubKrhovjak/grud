FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

COPY --from=builder --chown=appuser:appuser /app/server .

USER appuser

EXPOSE 8080

CMD ["./server"]
