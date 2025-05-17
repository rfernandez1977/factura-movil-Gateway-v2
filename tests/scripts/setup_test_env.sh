#!/bin/bash

# Script de configuración del ambiente de pruebas para FMgo

echo "🚀 Configurando ambiente de pruebas FMgo..."

# Verificar requisitos
check_requirements() {
    echo "📋 Verificando requisitos..."
    
    # Verificar Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go no está instalado"
        exit 1
    fi
    echo "✅ Go instalado: $(go version)"
    
    # Verificar Redis
    if ! command -v redis-cli &> /dev/null; then
        echo "❌ Redis no está instalado"
        exit 1
    fi
    echo "✅ Redis instalado: $(redis-cli --version)"
    
    # Verificar PostgreSQL
    if ! command -v psql &> /dev/null; then
        echo "❌ PostgreSQL no está instalado"
        exit 1
    fi
    echo "✅ PostgreSQL instalado: $(psql --version)"
}

# Configurar Redis para pruebas
setup_redis() {
    echo "🔄 Configurando Redis..."
    
    # Limpiar datos existentes
    redis-cli FLUSHALL
    
    # Configurar para pruebas
    redis-cli CONFIG SET maxmemory "1gb"
    redis-cli CONFIG SET maxmemory-policy "allkeys-lru"
    redis-cli CONFIG SET notify-keyspace-events "Ex"
    
    echo "✅ Redis configurado"
}

# Configurar base de datos de prueba
setup_database() {
    echo "🔄 Configurando base de datos..."
    
    # Variables de conexión
    DB_NAME="fmgo_test"
    DB_USER="postgres"
    
    # Crear base de datos de prueba
    psql -U $DB_USER -c "DROP DATABASE IF EXISTS $DB_NAME;"
    psql -U $DB_USER -c "CREATE DATABASE $DB_NAME;"
    
    # Aplicar migraciones
    echo "🔄 Aplicando migraciones..."
    go run cmd/migrate/main.go -env test
    
    echo "✅ Base de datos configurada"
}

# Generar certificados de prueba
generate_test_certs() {
    echo "🔄 Generando certificados de prueba..."
    
    # Crear directorio para certificados
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
    
    echo "✅ Certificados generados"
}

# Preparar datos de prueba
setup_test_data() {
    echo "🔄 Preparando datos de prueba..."
    
    # Copiar CAFs de prueba
    mkdir -p tests/data/caf
    cp test_cases/caf/*.xml tests/data/caf/
    
    # Copiar XMLs de ejemplo
    mkdir -p tests/data/xml
    cp test_cases/xml/*.xml tests/data/xml/
    
    echo "✅ Datos de prueba preparados"
}

# Configurar ambiente de pruebas
setup_test_environment() {
    echo "🔄 Configurando variables de ambiente..."
    
    # Crear archivo de configuración de pruebas
    cat > tests/config/test.env << EOF
FMGO_ENV=test
FMGO_DB_HOST=localhost
FMGO_DB_PORT=5432
FMGO_DB_NAME=fmgo_test
FMGO_DB_USER=postgres
FMGO_REDIS_HOST=localhost
FMGO_REDIS_PORT=6379
FMGO_SII_URL=https://maullin.sii.cl/DTEWS/
FMGO_CERT_PATH=tests/certs/test_client.crt
FMGO_KEY_PATH=tests/certs/test_client.key
EOF

    echo "✅ Variables de ambiente configuradas"
}

# Función principal
main() {
    echo "🚀 Iniciando configuración del ambiente de pruebas..."
    
    check_requirements
    setup_redis
    setup_database
    generate_test_certs
    setup_test_data
    setup_test_environment
    
    echo "✅ Ambiente de pruebas configurado exitosamente"
}

# Ejecutar script
main 