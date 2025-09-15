# Usa una imagen base de Go con Alpine para un contenedor ligero
FROM golang:1.22-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de dependencias y descarga los módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia el código fuente al directorio de trabajo
COPY . .

# Compila la aplicación
RUN go build -o search-api main.go

# Expone el puerto en el que corre la aplicación
EXPOSE 8082

# Comando para ejecutar la aplicación
CMD ["./search-api"]
