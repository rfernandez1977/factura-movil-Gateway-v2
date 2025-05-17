#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Verificar requisitos
check_requirements() {
    echo -e "${YELLOW}📋 Verificando requisitos...${NC}"
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go no está instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Go instalado: $(go version)${NC}"
    
    # Verificar Redis
    if ! command -v redis-cli &> /dev/null; then
        echo -e "${RED}❌ Redis no está instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Redis instalado: $(redis-cli --version)${NC}"
    
    # Verificar PostgreSQL
    if ! command -v psql &> /dev/null; then
        echo -e "${RED}❌ PostgreSQL no está instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ PostgreSQL instalado: $(psql --version)${NC}"
    
    # Verificar k6
    K6_PATH="/usr/local/Cellar/k6/1.0.0/bin/k6"
    if [ ! -f "$K6_PATH" ]; then
        echo -e "${RED}❌ k6 no está instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ k6 instalado: $($K6_PATH version)${NC}"
}

# Configurar Redis
setup_redis() {
    echo -e "\n${YELLOW}🔄 Configurando Redis...${NC}"
    
    # Limpiar datos existentes
    redis-cli FLUSHALL
    
    # Configurar para pruebas
    redis-cli CONFIG SET maxmemory "1gb"
    redis-cli CONFIG SET maxmemory-policy "allkeys-lru"
    redis-cli CONFIG SET notify-keyspace-events "Ex"
    
    echo -e "${GREEN}✅ Redis configurado${NC}"
}

# Configurar base de datos
setup_database() {
    echo -e "\n${YELLOW}🔄 Configurando base de datos...${NC}"
    
    # Variables de conexión
    DB_NAME="fmgo_test"
    DB_USER="postgres"
    DB_PASSWORD="test_password"
    
    export PGPASSWORD="${DB_PASSWORD}"
    
    # Crear base de datos de prueba
    psql -h localhost -U ${DB_USER} -d postgres -c "DROP DATABASE IF EXISTS ${DB_NAME};"
    psql -h localhost -U ${DB_USER} -d postgres -c "CREATE DATABASE ${DB_NAME};"
    
    # Aplicar migraciones si existen
    if [ -f "cmd/migrate/main.go" ]; then
        echo -e "${YELLOW}🔄 Aplicando migraciones...${NC}"
        go run cmd/migrate/main.go -env test
    fi
    
    echo -e "${GREEN}✅ Base de datos configurada${NC}"
}

# Generar certificados
generate_test_certs() {
    echo -e "\n${YELLOW}🔄 Generando certificados de prueba...${NC}"
    
    # Crear directorios
    mkdir -p tests/certs
    
    # Generar certificado raíz
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout tests/certs/test_ca.key \
        -out tests/certs/test_ca.crt \
        -subj "/CN=FMgo Test CA"
    
    # Generar certificado de cliente
    openssl req -newkey rsa:2048 -nodes \
        -keyout tests/certs/test_client.key \
        -out tests/certs/test_client.csr \
        -subj "/CN=FMgo Test Client"
    
    openssl x509 -req -days 365 \
        -in tests/certs/test_client.csr \
        -CA tests/certs/test_ca.crt \
        -CAkey tests/certs/test_ca.key \
        -CAcreateserial \
        -out tests/certs/test_client.crt
    
    echo -e "${GREEN}✅ Certificados generados${NC}"
}

# Preparar datos de prueba
setup_test_data() {
    echo -e "\n${YELLOW}🔄 Preparando datos de prueba...${NC}"
    
    # Crear directorios
    mkdir -p tests/data/{caf,xml}
    
    # Copiar datos de prueba si existen
    if [ -d "test_cases" ]; then
        cp -r test_cases/caf/*.xml tests/data/caf/ 2>/dev/null || true
        cp -r test_cases/xml/*.xml tests/data/xml/ 2>/dev/null || true
    fi
    
    if [ -d "examples" ]; then
        mkdir -p tests/data/json
        cp examples/test_*.json tests/data/json/ 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✅ Datos de prueba preparados${NC}"
}

# Configurar ambiente
setup_environment() {
    echo -e "\n${YELLOW}🔄 Configurando variables de ambiente...${NC}"
    
    # Crear archivo de configuración
    cat > tests/config/test.env << EOF
FMGO_ENV=test
FMGO_DB_HOST=localhost
FMGO_DB_PORT=5432
FMGO_DB_NAME=fmgo_test
FMGO_DB_USER=postgres
FMGO_DB_PASSWORD=test_password
FMGO_REDIS_HOST=localhost
FMGO_REDIS_PORT=6379
FMGO_SII_URL=https://maullin.sii.cl/DTEWS/
FMGO_CERT_PATH=tests/certs/test_client.crt
FMGO_KEY_PATH=tests/certs/test_client.key
EOF

    # Crear symlink para compatibilidad
    ln -sf tests/config/test.env .env.test
    
    echo -e "${GREEN}✅ Variables de ambiente configuradas${NC}"
}

# Función principal
main() {
    echo -e "${YELLOW}🚀 Iniciando configuración del ambiente de pruebas...${NC}"
    
    check_requirements
    setup_redis
    setup_database
    generate_test_certs
    setup_test_data
    setup_environment
    
    echo -e "\n${GREEN}✅ Ambiente de pruebas configurado exitosamente${NC}"
}

# Ejecutar script
main 