FROM golang:1.24.1 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates make curl
WORKDIR /root/

# Install migrate binary langsung di alpine
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin/migrate
RUN chmod +x /usr/local/bin/migrate

COPY --from=builder /app/app ./app
COPY --from=builder /app/db ./db

# Test migrate
RUN /usr/local/bin/migrate --version

ENV APP_ENV=production
ENV DB_HOST=postgres
EXPOSE 3000
CMD ["sh", "-c", "migrate -database \"$DATABASE_URL\" -path db/migrations up && ./app"]