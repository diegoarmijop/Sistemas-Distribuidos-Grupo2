# Usa una imagen base oficial de Go
FROM golang:1.23.4-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el módulo y descarga las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila la aplicación
RUN go build -o main .

# Expone el puerto en el que la aplicación se ejecuta
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]