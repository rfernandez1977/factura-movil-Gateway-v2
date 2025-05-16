#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuración
REPORT_DIR="test-reports/integration"
TIMEOUT="15m"
DOCKER_COMPOSE="dev/docker-compose.yml"

# Crear directorio de reportes
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}=== Ejecutando Tests de Integración ====${NC}\n"

# Función para mostrar el tiempo transcurrido
start_time=$(date +%s)
show_time_elapsed() {
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    echo -e "\n${BLUE}Tiempo total: ${elapsed}s${NC}"
}

# Trap para asegurar que siempre se muestre el tiempo
trap show_time_elapsed EXIT

# Verificar docker-compose
if [ ! -f "$DOCKER_COMPOSE" ]; then
    echo -e "${RED}Error: No se encuentra el archivo docker-compose.yml${NC}"
    exit 1
fi

# Iniciar servicios de prueba
echo -e "${YELLOW}Iniciando servicios de prueba...${NC}"
docker-compose -f "$DOCKER_COMPOSE" up -d

# Esperar a que los servicios estén listos
echo -e "${YELLOW}Esperando a que los servicios estén listos...${NC}"
sleep 10

# Función para verificar servicio
check_service() {
    local service=$1
    local port=$2
    echo -e "${YELLOW}Verificando $service...${NC}"
    for i in {1..30}; do
        if nc -z localhost $port; then
            echo -e "${GREEN}✓ $service está listo${NC}"
            return 0
        fi
        sleep 1
    done
    echo -e "${RED}✗ $service no está disponible${NC}"
    return 1
}

# Verificar servicios
check_service "PostgreSQL" 5432 && \
check_service "Redis" 6379 && \
check_service "MongoDB" 27017 && \
check_service "RabbitMQ" 5672

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: No todos los servicios están disponibles${NC}"
    docker-compose -f "$DOCKER_COMPOSE" down
    exit 1
fi

# Ejecutar tests de integración
echo -e "\n${YELLOW}Ejecutando tests de integración...${NC}"
gotestsum --format pkgname \
    --junitfile "$REPORT_DIR/junit.xml" \
    --jsonfile "$REPORT_DIR/tests.json" \
    -- -tags=integration \
    -timeout "$TIMEOUT" \
    ./...

test_exit_code=$?

# Detener servicios
echo -e "\n${YELLOW}Deteniendo servicios...${NC}"
docker-compose -f "$DOCKER_COMPOSE" down

# Mostrar resumen
echo -e "\n${BLUE}=== Resultados ===${NC}"
if [ $test_exit_code -eq 0 ]; then
    echo -e "Tests de Integración: ${GREEN}✓${NC}"
else
    echo -e "Tests de Integración: ${RED}✗${NC}"
fi

# Mostrar ubicación de reportes
echo -e "\n${BLUE}=== Reportes Generados ===${NC}"
echo -e "- Reporte JUnit: ${YELLOW}$REPORT_DIR/junit.xml${NC}"
echo -e "- Reporte JSON: ${YELLOW}$REPORT_DIR/tests.json${NC}"

exit $test_exit_code 