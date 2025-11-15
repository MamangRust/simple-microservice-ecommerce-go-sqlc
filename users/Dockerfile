FROM golang:1.25.4-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o user ./cmd/app/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate ./cmd/migrate/main.go

FROM alpine:latest


RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /var/log/app && \
    mkdir -p /app/pkg/database/migrations && \
    chown -R appuser:appuser /var/log/app /app

COPY --from=builder --chown=appuser:appuser /app/user .
COPY --from=builder --chown=appuser:appuser /app/migrate .


COPY --from=builder --chown=appuser:appuser /app/pkg/database/migrations/*.sql /app/pkg/database/migrations/

COPY --chown=appuser:appuser ./entrypoint.sh .

RUN chmod +x ./entrypoint.sh ./user ./migrate
USER appuser

EXPOSE 50052 8082

ENTRYPOINT ["./entrypoint.sh"]