# Documento Técnico: Gateway-XML

## Introducción
**Gateway-XML** es una aplicación escrita en Go que actúa como un intermediario entre sistemas internos y la API de Factura Móvil, permitiendo la gestión de documentos electrónicos (facturas, boletas, notas de crédito/débito, guías de despacho), clientes y productos. Este sistema está diseñado para ser modular, escalable y robusto, con integración de una base de datos PostgreSQL, métricas con Prometheus, y autenticación básica mediante API keys.

Este documento proporciona una guía completa para la instalación, configuración y desarrollo continuo del sistema.

---

## 1. Arquitectura y Componentes

### 1.1. Componentes Principales
- **Backend**:
  - Lenguaje: Go 1.21
  - Framework: Gin (para la API REST)
  - Base de Datos: PostgreSQL (almacenamiento local de documentos)
  - Métricas: Prometheus (monitoreo de solicitudes)
  - Logging: Biblioteca estándar de Go (`log`)

- **Integraciones**:
  - **Factura Móvil API**: Para la creación y consulta de documentos electrónicos, clientes y productos.
  - **Redis** (opcional): Configuración incluida para caché o colas, aunque no está implementada en el código actual.
  - **Grafana**: Visualización de métricas mediante un dashboard preconfigurado.

### 1.2. Estructura de Directorios
```
gateway/
├── main.go                  # Punto de entrada principal
├── handlers/                # Manejo de endpoints y lógica de negocio
│   ├── documents.go         # Handlers para documentos
│   └── entities.go          # Handlers para entidades
├── middleware/              # Middlewares (autenticación)
│   └── auth.go
├── db/                      # Conexión y operaciones con la base de datos
│   └── db.go
├── metrics/                 # Configuración de métricas con Prometheus
│   └── metrics.go
├── api/                     # Cliente HTTP para Factura Móvil
│   └── facturamovil.go
├── models/                  # Estructuras de datos
│   └── models.go
├── docs/                    # Documentación
│   └── gateway-api-endpoints.md
├── config/                  # Archivos de configuración
│   ├── redis-config.yaml    # Configuración de Redis
│   ├── workers-config.json  # Configuración de workers
│   └── grafana-dashboard.json  # Dashboard de Grafana
└── go.mod                   # Módulo de Go y dependencias
```

---

## 2. Requisitos del Sistema

### 2.1. Requisitos de Infraestructura
- **Sistema Operativo**: Linux (recomendado), macOS o Windows.
- **Go**: Versión 1.21 o superior.
- **PostgreSQL**: Versión 12 o superior.
- **Prometheus**: Para monitoreo de métricas.
- **Grafana**: Para visualización de métricas (opcional).
- **Redis**: Opcional, para caché o colas (no implementado en el código actual).

### 2.2. Dependencias de Software
Las dependencias están definidas en `go.mod`:
- `github.com/gin-gonic/gin v1.9.1`
- `github.com/lib/pq v1.10.9`
- `github.com/prometheus/client_golang v1.17.0`

Instala las dependencias ejecutando:
```bash
go mod tidy
```

---

## 3. Instalación

### 3.1. Descarga del Código
1. Clona el repositorio o copia los archivos proporcionados:
   - `main.go`
   - `handlers/documents.go`
   - `handlers/entities.go`
   - `middleware/auth.go`
   - `db/db.go`
   - `metrics/metrics.go`
   - `api/facturamovil.go`
   - `models/models.go`
   - `docs/gateway-api-endpoints.md`
   - `config/redis-config.yaml`
   - `config/workers-config.json`
   - `config/grafana-dashboard.json`
   - `go.mod`

2. Organiza los archivos en la estructura de directorios indicada.

### 3.2. Configuración de la Base de Datos (PostgreSQL)
1. Instala PostgreSQL si no está instalado:
   ```bash
   # Ejemplo para Ubuntu
   sudo apt update
   sudo apt install postgresql postgresql-contrib
   ```
2. Inicia el servicio de PostgreSQL:
   ```bash
   sudo service postgresql start
   ```
3. Crea una base de datos llamada `gateway`:
   ```bash
   sudo -u postgres psql -c "CREATE DATABASE gateway;"
   ```
4. Configura la cadena de conexión en `db/db.go`. Por defecto:
   ```go
   connStr := "host=localhost user=postgres password=secret dbname=gateway sslmode=disable"
   ```
   Ajusta `user`, `password`, y `host` según tu entorno.

### 3.3. Configuración de Prometheus
1. Descarga e instala Prometheus:
   ```bash
   wget https://github.com/prometheus/prometheus/releases/download/v2.47.0/prometheus-2.47.0.linux-amd64.tar.gz
   tar xvfz prometheus-2.47.0.linux-amd64.tar.gz
   cd prometheus-2.47.0.linux-amd64
   ```
2. Configura `prometheus.yml` para scrapear métricas del Gateway:
   ```yaml
   scrape_configs:
     - job_name: 'gateway'
       static_configs:
         - targets: ['localhost:3000']
   ```
3. Inicia Prometheus:
   ```bash
   ./prometheus --config.file=prometheus.yml
   ```

### 3.4. Configuración de Grafana (Opcional)
1. Instala Grafana:
   ```bash
   sudo apt-get install -y adduser libfontconfig1 musl
   wget https://dl.grafana.com/oss/release/grafana_10.1.0_amd64.deb
   sudo dpkg -i grafana_10.1.0_amd64.deb
   sudo systemctl start grafana-server
   ```
2. Accede a Grafana en `http://localhost:3000` (usuario: `admin`, contraseña: `admin`).
3. Agrega Prometheus como fuente de datos (URL: `http://localhost:9090`).
4. Importa el dashboard desde `config/grafana-dashboard.json`.

### 3.5. Configuración de Redis (Opcional)
1. Instala Redis:
   ```bash
   sudo apt install redis-server
   ```
2. Ajusta `config/redis-config.yaml` según tu entorno:
   ```yaml
   redis:
     host: "localhost"
     port: 6379
     password: ""
     db: 0
     max_retries: 3
     pool_size: 10
     min_idle_conns: 5
   ```
3. Inicia Redis:
   ```bash
   sudo systemctl start redis
   ```

### 3.6. Configuración de Workers (Opcional)
El archivo `config/workers-config.json` define la configuración de workers para tareas en segundo plano (no implementadas en el código actual):
```json
{
  "workers": {
    "document_processor": {
      "count": 5,
      "queue": "documents:pending",
      "retry_attempts": 3,
      "retry_delay_seconds": 10
    },
    "status_checker": {
      "count": 2,
      "queue": "status:pending",
      "retry_attempts": 2,
      "retry_delay_seconds": 5
    }
  }
}
```
Si deseas implementar workers, puedes usar una librería como `go-workers` y Redis como backend.

### 3.7. Configuración de Autenticación
1. En `middleware/auth.go`, reemplaza el API key por uno seguro:
   ```go
   if apiKey != "your-secure-api-key" {
   ```
2. Pasa el header `X-API-Key` en todas las solicitudes.

### 3.8. Ejecución del Servidor
1. Desde el directorio raíz del proyecto, ejecuta:
   ```bash
   go run main.go
   ```
2. El servidor estará disponible en `http://localhost:3000`.

---

## 4. Parametrización

### 4.1. Base de Datos
- **Cadena de Conexión**: Ajusta `connStr` en `db/db.go`.
- **Esquema**: La tabla `documents` se crea automáticamente al iniciar el servidor:
  ```sql
  CREATE TABLE documents (
      id SERIAL PRIMARY KEY,
      type VARCHAR(50) NOT NULL,
      data TEXT NOT NULL,
      created_at TIMESTAMP NOT NULL
  );
  ```

### 4.2. Factura Móvil API
- **URL Base y Token**: Definidos en `api/facturamovil.go`:
  ```go
  const (
      facturaMovilBaseURL = "http://produccion.facturamovil.cl"
      facmovToken         = "da395d31-7f91-424b-8034-cda17ab4ed83"
  )
  ```
- **Company ID**: Hardcoded como `29`. Ajusta según sea necesario.

### 4.3. Métricas
- Métricas disponibles en `/metrics`.
- Contador principal: `gateway_requests_total{endpoint, status}`.

### 4.4. Logging
- Los logs se escriben en la consola usando `log.Printf`.
- Ejemplo: `log.Printf("Failed to communicate with Factura Móvil for %s: %v", docType, err)`.

---

## 5. Funcionalidades

### 5.1. Endpoints Disponibles
Consulta `docs/gateway-api-endpoints.md` para una lista completa. Resumen:
- **Documentos**:
  - `POST /facturas`: Crea facturas.
  - `POST /boletas`: Crea boletas.
  - `POST /notas`: Crea notas de crédito/débito.
  - `POST /guias`: Crea guías de despacho.
- **Entidades**:
  - `POST /clientes`: Crea clientes.
  - `POST /productos`: Crea productos.
- **Consultas**:
  - `GET /documents/:id`: Consulta el estado de un documento.
  - `GET /documents/:id/pdf`: Descarga el PDF de un documento.
- **Métricas**:
  - `GET /metrics`: Métricas de Prometheus.

### 5.2. Características
- **Almacenamiento Local**: Documentos se guardan en PostgreSQL.
- **Validaciones**: Campos obligatorios (`date`, `details`, `netTotal` para documentos; `code`, `name`, `address` para clientes).
- **Reintentos**: Hasta 3 reintentos para solicitudes a Factura Móvil.
- **Autenticación**: API key requerida en el header `X-API-Key`.

---

## 6. Guía para Desarrolladores

### 6.1. Estructura del Código
- **Handlers** (`handlers/`): Lógica de los endpoints.
- **Middleware** (`middleware/`): Autenticación y otros middlewares.
- **Base de Datos** (`db/`): Conexión y operaciones con PostgreSQL.
- **API** (`api/`): Cliente HTTP para Factura Móvil.
- **Modelos** (`models/`): Estructuras de datos.
- **Métricas** (`metrics/`): Configuración de Prometheus.

### 6.2. Agregar Nuevos Endpoints
1. Define el handler en `handlers/`. Ejemplo para un nuevo endpoint `POST /new-endpoint`:
   ```go
   // handlers/new.go
   package handlers

   import "github.com/gin-gonic/gin"

   func NewEndpointHandler() gin.HandlerFunc {
       return func(c *gin.Context) {
           c.JSON(http.StatusOK, gin.H{"status": "new endpoint"})
       }
   }
   ```
2. Registra el endpoint en `main.go`:
   ```go
   r.POST("/new-endpoint", handlers.NewEndpointHandler())
   ```

### 6.3. Extender la Base de Datos
1. Añade nuevas tablas o columnas en `db/db.go`.
2. Actualiza los modelos en `models/models.go` si es necesario.

### 6.4. Integrar Redis
1. Usa `config/redis-config.yaml` para configurar Redis.
2. Integra una librería como `github.com/go-redis/redis` para conectar y usar Redis.

### 6.5. Mejorar el Monitoreo
- Agrega nuevos contadores o histogramas en `metrics/metrics.go`.
- Actualiza el dashboard de Grafana en `config/grafana-dashboard.json`.

---

## 7. Solución de Problemas

### 7.1. Errores Comunes
- **No se puede conectar a PostgreSQL**:
  - Verifica la cadena de conexión en `db/db.go`.
  - Asegúrate de que PostgreSQL esté corriendo (`sudo systemctl status postgresql`).
- **Error 401 (Unauthorized)**:
  - Asegúrate de pasar un `X-API-Key` válido en el header.
- **Factura Móvil API falla**:
  - Verifica la URL base y el token en `api/facturamovil.go`.
  - Revisa los logs para detalles del error.
- **Métricas no visibles en Prometheus**:
  - Asegúrate de que Prometheus esté scrapeando `http://localhost:3000/metrics`.

### 7.2. Logs
- Los logs se escriben en la consola.
- Busca mensajes como `Failed to communicate with Factura Móvil` para errores específicos.

---

## 8. Consideraciones de Escalabilidad

- **Base de Datos**: Usa índices en la tabla `documents` para mejorar el rendimiento:
  ```sql
  CREATE INDEX idx_documents_type ON documents(type);
  ```
- **Caché**: Implementa Redis para cachear respuestas frecuentes.
- **Workers**: Usa `config/workers-config.json` para procesar tareas en segundo plano (por ejemplo, reintentos de Factura Móvil).
- **Balanceo de Carga**: Despliega múltiples instancias del Gateway detrás de un balanceador de carga como Nginx.

---

## 9. Contacto y Soporte
Para soporte adicional, contacta al equipo de desarrollo o consulta la documentación en `docs/`.

**Última Actualización**: 24 de abril de 2025