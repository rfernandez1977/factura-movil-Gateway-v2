#!/bin/bash

# Script para ejecutar pruebas de integración
# Este script ejecuta las pruebas de integración con el SII

# Colores para mensajes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Función para imprimir mensajes de estado
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Verificar que estamos en el directorio correcto
if [ ! -d "test/config" ]; then
    print_error "Este script debe ejecutarse desde el directorio raíz del proyecto"
    exit 1
fi

# Verificar que los datos de prueba existen
if [ ! -d "test/data" ]; then
    print_error "No se encontraron los datos de prueba"
    echo "Por favor, ejecute primero los scripts de generación de datos"
    exit 1
fi

# Ejecutar pruebas de autenticación
echo "Ejecutando pruebas de autenticación..."
go test -v ./test/integration/auth/... -tags=integration

# Ejecutar pruebas de firma
echo "Ejecutando pruebas de firma..."
go test -v ./test/integration/firma/... -tags=integration

# Ejecutar pruebas de documentos
echo "Ejecutando pruebas de documentos..."
go test -v ./test/integration/dte/... -tags=integration

# Ejecutar pruebas de CAF
echo "Ejecutando pruebas de CAF..."
go test -v ./test/integration/caf/... -tags=integration

# Ejecutar pruebas de certificados
echo "Ejecutando pruebas de certificados..."
go test -v ./test/integration/certs/... -tags=integration

# Ejecutar pruebas de envío
echo "Ejecutando pruebas de envío..."
go test -v ./test/integration/envio/... -tags=integration

# Ejecutar pruebas de consulta
echo "Ejecutando pruebas de consulta..."
go test -v ./test/integration/consulta/... -tags=integration

# Generar reporte de cobertura
echo "Generando reporte de cobertura..."
go test -coverprofile=coverage.out ./test/integration/...
go tool cover -html=coverage.out -o coverage.html

print_status "Pruebas de integración completadas"
echo "Por favor, revise el reporte de cobertura en coverage.html" 