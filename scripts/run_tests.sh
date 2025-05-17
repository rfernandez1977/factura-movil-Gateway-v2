#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Configuración
COVERAGE_THRESHOLD=80
TEST_TIMEOUT="10m"
PARALLEL_TESTS=4

# Verificar ambiente
if [ ! -f "tests/config/test.env" ]; then
    echo -e "${RED}❌ Ambiente de pruebas no configurado. Ejecute setup_test_env.sh primero${NC}"
    exit 1
fi

# Cargar variables de ambiente
source tests/config/test.env

# Función para ejecutar pruebas unitarias
run_unit_tests() {
    echo -e "${YELLOW}📋 Ejecutando pruebas unitarias...${NC}"
    
    # Ejecutar pruebas con coverage
    go test -v -timeout ${TEST_TIMEOUT} -coverprofile=coverage.out \
        -covermode=atomic -parallel ${PARALLEL_TESTS} \
        ./... -tags=unit
    
    # Verificar cobertura
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    
    if (( $(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
        echo -e "${RED}❌ Cobertura insuficiente: ${COVERAGE}% (mínimo: ${COVERAGE_THRESHOLD}%)${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Cobertura: ${COVERAGE}%${NC}"
    fi
    
    # Generar reporte HTML
    go tool cover -html=coverage.out -o coverage.html
}

# Función para ejecutar pruebas de integración
run_integration_tests() {
    echo -e "${YELLOW}🔄 Ejecutando pruebas de integración...${NC}"
    
    go test -v -timeout ${TEST_TIMEOUT} \
        -tags=integration \
        ./tests/integration/...
}

# Función para ejecutar pruebas de carga
run_load_tests() {
    echo -e "${YELLOW}🔨 Ejecutando pruebas de carga...${NC}"
    
    # Escenario: Carga Normal
    echo -e "${YELLOW}📊 Escenario: Carga Normal (100 RPS)${NC}"
    k6 run tests/load/normal_load.js
    
    # Escenario: Pico de Carga
    echo -e "${YELLOW}📊 Escenario: Pico de Carga (500 RPS)${NC}"
    k6 run tests/load/peak_load.js
    
    # Escenario: Prueba de Estrés
    echo -e "${YELLOW}📊 Escenario: Prueba de Estrés${NC}"
    k6 run tests/load/stress_test.js
}

# Función para generar reporte
generate_report() {
    echo -e "${YELLOW}📝 Generando reporte de pruebas...${NC}"
    
    # Crear directorio de reportes
    REPORT_DIR="tests/reports/$(date +%Y%m%d_%H%M%S)"
    mkdir -p ${REPORT_DIR}
    
    # Copiar resultados
    cp coverage.html ${REPORT_DIR}/
    cp coverage.out ${REPORT_DIR}/
    
    # Generar reporte resumen
    cat > ${REPORT_DIR}/summary.md << EOF
# Reporte de Pruebas FMgo

## Resumen
- Fecha: $(date)
- Versión: $(git describe --tags --always)
- Commit: $(git rev-parse HEAD)
- Duración total: ${SECONDS}s

## Resultados
- Pruebas unitarias: ${UNIT_TESTS_RESULT}
- Pruebas de integración: ${INTEGRATION_TESTS_RESULT}
- Pruebas de carga: ${LOAD_TESTS_RESULT}

## Métricas
- Cobertura: ${COVERAGE}%
- Latencia P95: $(jq '.metrics.http_req_duration.p95' k6_report.json)ms
- Throughput: $(jq '.metrics.http_reqs.rate' k6_report.json) req/s
- Error rate: $(jq '.metrics.errors' k6_report.json)%

## Detalles
Ver archivos adjuntos para más información.
EOF

    echo -e "${GREEN}✅ Reporte generado en ${REPORT_DIR}${NC}"
}

# Función principal
main() {
    echo -e "${YELLOW}🚀 Iniciando suite de pruebas FMgo...${NC}"
    SECONDS=0
    
    # Ejecutar pruebas
    if run_unit_tests; then
        UNIT_TESTS_RESULT="✅ PASS"
    else
        UNIT_TESTS_RESULT="❌ FAIL"
        echo -e "${RED}❌ Pruebas unitarias fallidas${NC}"
        exit 1
    fi
    
    if run_integration_tests; then
        INTEGRATION_TESTS_RESULT="✅ PASS"
    else
        INTEGRATION_TESTS_RESULT="❌ FAIL"
        echo -e "${RED}❌ Pruebas de integración fallidas${NC}"
        exit 1
    fi
    
    if run_load_tests; then
        LOAD_TESTS_RESULT="✅ PASS"
    else
        LOAD_TESTS_RESULT="❌ FAIL"
        echo -e "${RED}❌ Pruebas de carga fallidas${NC}"
        exit 1
    fi
    
    # Generar reporte
    generate_report
    
    echo -e "${GREEN}✅ Suite de pruebas completada exitosamente${NC}"
    echo "⏱️ Tiempo total: ${SECONDS}s"
}

# Ejecutar script
main 