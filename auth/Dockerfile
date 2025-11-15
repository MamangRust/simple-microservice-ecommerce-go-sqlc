FROM golang:1.25.4-alpine3.22 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o auth ./cmd/app/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate ./cmd/migrate/main.go

FROM alpine:3.22

RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /var/log/app && \
    mkdir -p /app/pkg/database/migrations && \
    chown -R appuser:appuser /var/log/app /app

COPY --from=builder --chown=appuser:appuser /app/auth .
COPY --from=builder --chown=appuser:appuser /app/migrate .


COPY --from=builder --chown=appuser:appuser /app/pkg/database/migrations/*.sql /app/pkg/database/migrations/

COPY --chown=appuser:appuser ./entrypoint.sh .

RUN chmod +x ./entrypoint.sh ./auth ./migrate

USER appuser

EXPOSE 50051 8081

ENTRYPOINT ["./entrypoint.sh"]