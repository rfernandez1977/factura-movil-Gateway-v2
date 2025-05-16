#!/bin/bash

# Script de inicialización del ambiente de desarrollo

echo "Iniciando configuración del ambiente de desarrollo..."

# Crear directorios necesarios
mkdir -p {logs,temp,certs}

# Verificar dependencias
echo "Verificando dependencias..."
command -v go >/dev/null 2>&1 || { echo "Go no está instalado. Por favor, instale Go 1.21 o superior."; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Docker no está instalado. Por favor, instale Docker."; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "Docker Compose no está instalado."; exit 1; }

# Configurar variables de entorno
echo "Configurando variables de entorno..."
if [ ! -f .env ]; then
    cat > .env << EOF
# Variables de Base de Datos
DB_USER=fmgo_dev
DB_PASSWORD=development_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fmgo_dev

# Redis
REDIS_PASSWORD=development_password

# SII (Certificación)
SII_CERT_PATH=./certs/cert.pfx
SII_KEY_PATH=./certs/key.pem

# Seguridad
JWT_SECRET=development_secret_key

# Ambiente
ENV=development
EOF
    echo "Archivo .env creado con valores por defecto"
else
    echo "Archivo .env ya existe"
fi

# Inicializar Go modules
echo "Inicializando módulos Go..."
go mod tidy

# Configurar Git hooks
echo "Configurando Git hooks..."
cat > .git/hooks/pre-commit << EOF
#!/bin/bash
go fmt ./...
go vet ./...
EOF
chmod +x .git/hooks/pre-commit

# Crear docker-compose para ambiente de desarrollo
echo "Creando docker-compose.yml..."
cat > docker-compose.yml << EOF
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    environment:
      POSTGRES_USER: \${DB_USER}
      POSTGRES_PASSWORD: \${DB_PASSWORD}
      POSTGRES_DB: \${DB_NAME}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:6-alpine
    command: redis-server --requirepass \${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  prometheus:
    image: prom/prometheus:v2.30.3
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  grafana:
    image: grafana/grafana:8.2.2
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  postgres_data:
  redis_data:
  prometheus_data:
  grafana_data:
EOF

# Crear configuración de Prometheus
mkdir -p config
cat > config/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'fmgo'
    static_configs:
      - targets: ['localhost:9090']
EOF

echo "Iniciando servicios Docker..."
docker-compose up -d

echo "Esperando que los servicios estén listos..."
sleep 10

echo "Configuración completada."
echo "Para comenzar a desarrollar:"
echo "1. Revise y ajuste el archivo .env según sea necesario"
echo "2. Los servicios de base de datos y caché están corriendo en Docker"
echo "3. Puede acceder a Grafana en http://localhost:3000"
echo "4. Los logs se almacenarán en el directorio ./logs" 