# Guía de Desarrollo

## Estructura del Proyecto

### Directorios Principales

- `api/`: Definiciones de API y documentación
- `config/`: Archivos de configuración
- `controllers/`: Controladores de la aplicación
- `db/`: Configuración y conexión a base de datos
- `docs/`: Documentación del proyecto
- `handlers/`: Manejadores de peticiones HTTP
- `middleware/`: Middleware de la aplicación
- `migrations/`: Migraciones de base de datos
- `models/`: Modelos de datos
- `repository/`: Repositorios para acceso a datos
- `routes/`: Definición de rutas
- `services/`: Servicios de negocio
- `utils/`: Utilidades generales
- `metrics/`: Métricas y monitoreo

## Configuración del Entorno de Desarrollo

### Requisitos

- Go 1.24 o superior
- MongoDB
- Grafana (para el dashboard)
- Certificados digitales para firma electrónica

### Pasos de Configuración

1. Clonar el repositorio:
```bash
git clone https://github.com/tu-usuario/fmgo.git
cd fmgo
```

2. Instalar dependencias:
```bash
go mod download
```

3. Configurar variables de entorno:
```bash
cp .env.example .env
# Editar .env con tus configuraciones
```

4. Iniciar MongoDB:
```bash
mongod --dbpath /ruta/a/tu/directorio/datos
```

5. Iniciar Grafana:
```bash
docker run -d -p 3000:3000 grafana/grafana
```

## Convenciones de Código

### Estructura de Archivos

- Los nombres de archivos deben ser en minúsculas con palabras separadas por guiones bajos
- Los archivos de prueba deben terminar en `_test.go`
- Los archivos de configuración deben estar en el directorio `config/`

### Convenciones de Nombrado

- Paquetes: nombres en minúsculas, una palabra
- Interfaces: nombres en PascalCase, terminando en "er" cuando sea apropiado
- Variables: camelCase
- Constantes: UPPER_CASE
- Funciones: PascalCase para funciones públicas, camelCase para privadas

### Documentación

- Cada paquete debe tener un archivo `doc.go` con documentación
- Las funciones públicas deben tener comentarios que expliquen su propósito
- Usar ejemplos de código cuando sea apropiado

## Desarrollo de Nuevas Características

### 1. Crear una Rama

```bash
git checkout -b feature/nombre-de-la-caracteristica
```

### 2. Implementar la Característica

- Crear los modelos necesarios en `models/`
- Implementar la lógica de negocio en `services/`
- Crear los controladores en `controllers/`
- Definir las rutas en `routes/`
- Agregar pruebas unitarias

### 3. Pruebas

```bash
go test ./...
```

### 4. Documentación

- Actualizar la documentación en `docs/`
- Agregar ejemplos de uso
- Documentar cambios en la API

### 5. Crear Pull Request

- Asegurarse de que todas las pruebas pasen
- Actualizar la documentación
- Solicitar revisión de código

## Manejo de Errores

### Estructura de Error

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### Códigos de Error

- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 429: Too Many Requests
- 500: Internal Server Error

## Logging

### Configuración

```go
logger, _ := zap.NewProduction()
defer logger.Sync()
```

### Niveles de Log

- DEBUG: Información detallada para debugging
- INFO: Información general de la aplicación
- WARN: Advertencias que no afectan la funcionalidad
- ERROR: Errores que afectan la funcionalidad
- FATAL: Errores críticos que detienen la aplicación

## Monitoreo

### Métricas

- Usar Prometheus para métricas
- Agregar métricas personalizadas cuando sea necesario
- Documentar nuevas métricas

### Dashboard

- Actualizar el dashboard de Grafana cuando se agreguen nuevas métricas
- Mantener la documentación del dashboard actualizada

## Seguridad

### Autenticación

- Usar API Keys para autenticación
- Implementar rate limiting
- Validar todas las entradas

### Certificados

- Mantener los certificados actualizados
- Usar certificados válidos para producción
- Rotar las claves regularmente

## Despliegue

### Requisitos

- Servidor con Go 1.24 o superior
- MongoDB
- Grafana
- Certificados digitales

### Pasos

1. Compilar la aplicación:
```bash
go build -o fmgo
```

2. Configurar el entorno:
```bash
export SII_BASE_URL=https://palena.sii.cl
export CERT_PATH=/ruta/a/cert.pem
export KEY_PATH=/ruta/a/key.pem
```

3. Iniciar la aplicación:
```bash
./fmgo
```

## Mantenimiento

### Actualizaciones

- Mantener las dependencias actualizadas
- Revisar y aplicar parches de seguridad
- Actualizar la documentación

### Monitoreo

- Revisar logs regularmente
- Monitorear métricas
- Responder a alertas

## Contribución

### Proceso

1. Fork el repositorio
2. Crear una rama para tu feature
3. Implementar los cambios
4. Agregar pruebas
5. Actualizar documentación
6. Crear Pull Request

### Código de Conducta

- Respetar a otros contribuidores
- Mantener discusiones constructivas
- Seguir las convenciones de código
- Mantener la documentación actualizada 