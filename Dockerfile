# Build stage
FROM golang:1.22.4-alpine3.20 AS builder
# Establecer el directorio de trabajo
WORKDIR /app
# Copiar todos los archivos de la raiz 
COPY . .
# Compile the app and add the binary to a `main` folder
RUN go build -o main main.go
# Install curl
RUN apk add curl
# Install to golang-migrate binary library
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]