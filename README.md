
# Sistema de Gestión de Cursos - Arquitectura de Microservicios

## Descripción

Sistema completo de gestión de cursos desarrollado con arquitectura de microservicios, que permite la creación, búsqueda, inscripción y administración de cursos educativos. El sistema incluye un frontend moderno en React y múltiples microservicios backend desarrollados en Go.

## Arquitectura del Sistema

### Microservicios

- **API Principal (Courses-API)**: Gestión CRUD de cursos con MongoDB
- **API de Usuarios (Users-API)**: Autenticación y gestión de usuarios con MySQL + Memcached
- **API de Búsqueda (Search-API)**: Búsqueda avanzada con SolR
- **API de Inscripciones (Inscriptions-API)**: Gestión de inscripciones con balanceo de carga
- **Frontend**: Interfaz de usuario en React con diseño moderno

### Infraestructura

- **Docker & Docker Compose**: Contenedorización completa
- **RabbitMQ**: Mensajería para sincronización entre servicios
- **SolR**: Motor de búsqueda avanzada
- **MongoDB**: Base de datos principal para cursos
- **MySQL**: Base de datos para usuarios
- **Memcached**: Caché para optimización de consultas
- **Nginx**: Balanceador de carga para inscripciones

## Características Principales

### Frontend
-  **Interfaz moderna y responsiva** con diseño Material Design
-  **Búsqueda en tiempo real** de cursos por nombre, descripción y categoría
-  **Gestión de sesiones** con JWT
-  **Vistas diferenciadas** para usuarios y administradores
-  **Subida de imágenes** personalizadas para cursos
-  **Navegación intuitiva** con React Router

### Backend
-  **Arquitectura de microservicios** con Go y Gin
-  **Patrón MVC** implementado en todos los servicios
-  **Autenticación JWT** con encriptación bcrypt
-  **Caché con Memcached** para optimización
-  **Búsqueda avanzada** con SolR
-  **Sincronización en tiempo real** con RabbitMQ
-  **Operaciones CRUD** completas en MongoDB
-  **Cálculo concurrente** de disponibilidad con Go Routines

## Tecnologías Utilizadas

### Frontend
- **React 18** - Biblioteca de interfaz de usuario
- **React Router DOM** - Enrutamiento
- **Axios** - Cliente HTTP
- **CSS3** - Estilos modernos y responsivos

### Backend
- **Go 1.21+** - Lenguaje de programación
- **Gin** - Framework web
- **GORM** - ORM para bases de datos
- **JWT-Go** - Autenticación con tokens
- **bcrypt** - Encriptación de contraseñas

### Bases de Datos
- **MongoDB** - Base de datos principal (cursos)
- **MySQL 8** - Base de datos de usuarios
- **SolR 8.11.1** - Motor de búsqueda
- **Memcached** - Caché en memoria

### Infraestructura
- **Docker & Docker Compose** - Contenedorización
- **RabbitMQ** - Mensajería asíncrona
- **Nginx** - Balanceador de carga
- **GitHub** - Control de versiones

## Instalación y Configuración

### Prerrequisitos
- Docker y Docker Compose instalados
- Git

### Pasos de Instalación

1. **Clonar el repositorio**
   ```bash
   git clone <https://github.com/martubecerra/final-arqsoft2.git>
   cd FinalArqsoft2
   ```

2. **Iniciar todos los servicios**
   ```bash
   docker-compose up -d
   ```

3. **Verificar que todos los servicios estén corriendo**
   ```bash
   docker-compose ps
   ```

4. **Acceder a la aplicación**
   - Frontend: http://localhost:3000
   - API Principal: http://localhost:8080
   - API Usuarios: http://localhost:8083
   - API Búsqueda: http://localhost:8082
   - API Inscripciones: http://localhost:8081
   - SolR Admin: http://localhost:8983
   - RabbitMQ Admin: http://localhost:15672

### Configuración de Servicios

#### Credenciales por defecto:
- **MySQL**: root/root
- **MongoDB**: root/root
- **RabbitMQ**: root/root
- **SolR**: Sin autenticación

## Uso del Sistema

### Para Usuarios Regulares
1. **Registro/Login**: Crear cuenta o iniciar sesión
2. **Explorar Cursos**: Ver cursos disponibles en el Home
3. **Buscar Cursos**: Usar la barra de búsqueda para encontrar cursos específicos
4. **Ver Detalles**: Hacer clic en un curso para ver información completa
5. **Inscribirse**: Registrarse en cursos de interés

### Para Administradores
1. **Gestión de Cursos**: Acceder a "Gestión de Cursos" desde el menú
2. **Crear Cursos**: Agregar nuevos cursos con imágenes personalizadas
3. **Editar Cursos**: Modificar información de cursos existentes
4. **Eliminar Cursos**: Remover cursos del sistema
5. **Ver Inscripciones**: Monitorear inscripciones de usuarios

## Endpoints Principales

### API de Cursos (Puerto 8080)
```
GET    /courses          - Obtener todos los cursos
POST   /courses          - Crear nuevo curso
GET    /courses/:id      - Obtener curso por ID
PUT    /courses/:id      - Actualizar curso
DELETE /courses/:id      - Eliminar curso
```

### API de Usuarios (Puerto 8083)
```
POST   /register         - Registrar usuario
POST   /login            - Iniciar sesión
GET    /users/:id        - Obtener usuario por ID
```

### API de Búsqueda (Puerto 8082)
```
GET    /search?q=query   - Buscar cursos
GET    /search/filter    - Filtrar por capacidad
```

### API de Inscripciones (Puerto 8081)
```
POST   /inscriptions     - Crear inscripción
GET    /inscriptions     - Obtener inscripciones
```

## Testing

### Verificación de Funcionalidades

1. **Conexión Frontend-Microservicios**
   - El frontend consume correctamente las APIs
   - Las peticiones HTTP funcionan sin errores

2. **Búsqueda y Filtrado**
   - Búsqueda por nombre, descripción y categoría
   - Filtrado por capacidad en SolR
   - Resultados relevantes y ordenados

3. **Gestión de Usuarios**
   - Registro e inicio de sesión
   - Autenticación JWT
   - Caché con Memcached

4. **Operaciones CRUD**
   - Creación, lectura, actualización y eliminación de cursos
   - Sincronización con RabbitMQ
   - Indexación automática en SolR

5. **Inscripciones**
   - Proceso de inscripción funcional
   - Balanceo de carga con Nginx
   - Cálculo concurrente de disponibilidad

## Monitoreo y Logs

### Verificar Estado de Servicios
```bash
# Ver logs de todos los servicios
docker-compose logs

# Ver logs de un servicio específico
docker-compose logs courses-api

# Ver estado de contenedores
docker-compose ps
```

## Solución de Problemas

### Problemas Comunes

1. **Puertos ocupados**
   ```bash
   # Verificar puertos en uso
   lsof -i :3000
   lsof -i :8080
   
   # Detener servicios que usen los puertos
   docker-compose down
   ```

2. **Problemas de memoria**
   ```bash
   # Limpiar recursos Docker
   docker system prune -a
   docker volume prune
   ```

3. **Base de datos vacía**
   - Los datos se inicializan automáticamente
   - Para resetear: `docker-compose down -v && docker-compose up -d`

4. **Problemas de sincronización SolR**
   - Verificar RabbitMQ: http://localhost:15672
   - Reindexar manualmente si es necesario

##  Estructura del Proyecto

```
FinalArqsoft2/
├── frontend/                 # Aplicación React
│   ├── src/
│   │   ├── components/       # Componentes React
│   │   ├── assets/          # Estilos y recursos
│   │   └── context/         # Contexto de usuario
│   ├── public/              # Archivos públicos
│   └── Dockerfile           # Contenedor del frontend
├── courses-api/             # Microservicio de cursos
│   ├── controllers/         # Controladores MVC
│   ├── services/           # Lógica de negocio
│   ├── repositories/       # Acceso a datos
│   ├── domain/            # Modelos de dominio
│   └── main.go            # Punto de entrada
├── users-api/              # Microservicio de usuarios
├── search-api/             # Microservicio de búsqueda
├── inscriptions-api/       # Microservicio de inscripciones
├── mysql-init/             # Scripts de inicialización MySQL
├── docker-compose.yml      # Orquestación de servicios
└── README.md              # Este archivo
```

## Autores

- **Martina Becerra** - *2214822* 
- **Sofía Contreras** - *2215803* 
- **Manuel Frias** - *2217890* 
- **José Ruarte** - *2206224* 

## Agradecimientos

- Arquitectura de Software II - Universidad Católica de Córdoba
- Profesores y compañeros del curso

---
