FROM golang:1.23.4-alpine

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY . . 

RUN go build -o main .

# Exponer el puerto 8080
EXPOSE 8080

# Copiar el script de entrada
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Usar el script como comando por defecto
CMD ["/entrypoint.sh"]
