#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuración
COVERAGE_DIR="coverage"
REPORT_DIR="test-reports"
MIN_COVERAGE=80
TIMEOUT="10m"

# Crear directorios si no existen
mkdir -p "$COVERAGE_DIR"
mkdir -p "$REPORT_DIR"

echo -e "${BLUE}=== Ejecutando Suite de Tests ====${NC}\n"

# Función para mostrar el tiempo transcurrido
start_time=$(date +%s)
show_time_elapsed() {
    end_time=$(date +%s)
    elapsed=$((end_time - start_time))
    echo -e "\n${BLUE}Tiempo total: ${elapsed}s${NC}"
}

# Trap para asegurar que siempre se muestre el tiempo
trap show_time_elapsed EXIT

# Limpiar reportes anteriores
echo -e "${YELLOW}Limpiando reportes anteriores...${NC}"
rm -f "$COVERAGE_DIR"/* "$REPORT_DIR"/*

# Ejecutar tests con gotestsum y generar reporte JUnit
echo -e "\n${YELLOW}Ejecutando tests unitarios...${NC}"
gotestsum --format pkgname \
    --junitfile "$REPORT_DIR/junit.xml" \
    --jsonfile "$REPORT_DIR/tests.json" \
    -- -timeout "$TIMEOUT" \
    -coverprofile="$COVERAGE_DIR/coverage.out" \
    -covermode=atomic \
    ./...

test_exit_code=$?

# Generar reporte de cobertura HTML
echo -e "\n${YELLOW}Generando reporte de cobertura...${NC}"
go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$COVERAGE_DIR/coverage.html"

# Calcular porcentaje de cobertura
coverage_percent=$(go tool cover -func="$COVERAGE_DIR/coverage.out" | grep total | awk '{print substr($3, 1, length($3)-1)}')

echo -e "\n${BLUE}=== Resumen de Tests ===${NC}"
echo -e "Cobertura total: ${YELLOW}${coverage_percent}%${NC}"

# Verificar cobertura mínima
coverage_status="✓"
if (( $(echo "$coverage_percent < $MIN_COVERAGE" | bc -l) )); then
    coverage_status="✗"
    echo -e "${RED}La cobertura está por debajo del mínimo requerido ($MIN_COVERAGE%)${NC}"
fi

# Mostrar resumen
echo -e "\n${BLUE}=== Resultados ===${NC}"
echo -e "Tests: ${test_exit_code == 0 ? "${GREEN}✓" : "${RED}✗"}${NC}"
echo -e "Cobertura: ${coverage_percent >= MIN_COVERAGE ? "${GREEN}$coverage_status" : "${RED}$coverage_status"}${NC} ($coverage_percent%)"

# Mostrar ubicación de reportes
echo -e "\n${BLUE}=== Reportes Generados ===${NC}"
echo -e "- Reporte JUnit: ${YELLOW}$REPORT_DIR/junit.xml${NC}"
echo -e "- Reporte JSON: ${YELLOW}$REPORT_DIR/tests.json${NC}"
echo -e "- Cobertura HTML: ${YELLOW}$COVERAGE_DIR/coverage.html${NC}"

# Salir con el código de error apropiado
if [ $test_exit_code -ne 0 ] || (( $(echo "$coverage_percent < $MIN_COVERAGE" | bc -l) )); then
    exit 1
fi

exit 0 