FROM golang:1.25-alpine AS builder

WORKDIR /app

#install app deps
COPY go.mod .
COPY go.sum .

#Download the Go module deps
RUN go mod download

COPY . .


#copiamos el binario
#esto generare un archivo ejecutable llamado server,binario
RUN go build -o server  ./cmd/server

#imagen mas ligera
FROM alpine:latest

WORKDIR /app
#copiamos el binario
COPY --from=builder /app/server .
# Instalar Air
RUN go install github.com/air-verse/air@latest

COPY .env ./

EXPOSE 8081
#

#comandos de arranque   ejecutmaos el binario creado en la fase de build que hemos llamado server -o server(este es el nombre del binaro)
CMD ["air", "-c", ".air.toml"]
 # CMD ["./server"]
