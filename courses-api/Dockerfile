# Usa la imagen oficial de Go como base
FROM golang:1.22-alpine

# Establece el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los archivos go.mod y go.sum
COPY go.mod go.sum ./

# Descarga todas las dependencias
RUN go mod download

# Copia el código fuente del proyecto al contenedor
COPY . .

# La carpeta de imágenes ya está en el contexto, no es necesario copiarla

# Compila la aplicación
RUN go build -o main .

# Expone el puerto 8080
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]

