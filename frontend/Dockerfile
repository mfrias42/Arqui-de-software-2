#Usa una imagen base de Node.js
FROM node:18

#Setea el directorio de trabajo
WORKDIR /app

#Copia los archivos de dependencias
COPY package*.json ./

#Instala las dependencias
RUN npm install

#Copia el resto del código del frontend
COPY . .

#Construye la aplicación de React
RUN npm run build

#Instala un servidor HTTP estático
RUN npm install -g serve

#Expone el puerto en el que el frontend se ejecutará
EXPOSE 3000

#Comando para ejecutar la aplicación
CMD ["serve", "-s", "build"]
