#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuración
COVERAGE_DIR="coverage"
REPORT_DIR="$COVERAGE_DIR/reports"
PROFILE_DIR="$COVERAGE_DIR/profiles"
MIN_COVERAGE=80

# Crear directorios
mkdir -p "$REPORT_DIR"
mkdir -p "$PROFILE_DIR"

echo -e "${BLUE}=== Generando Reportes de Cobertura ====${NC}\n"

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
rm -f "$COVERAGE_DIR"/*.out "$REPORT_DIR"/* "$PROFILE_DIR"/*

# Generar perfiles de cobertura por paquete
echo -e "\n${YELLOW}Generando perfiles de cobertura...${NC}"
packages=$(go list ./...)
for pkg in $packages; do
    pkg_name=$(basename $pkg)
    echo -e "Testing ${YELLOW}$pkg_name${NC}..."
    go test -coverprofile="$PROFILE_DIR/$pkg_name.out" "$pkg"
done

# Combinar perfiles de cobertura
echo -e "\n${YELLOW}Combinando perfiles de cobertura...${NC}"
echo "mode: atomic" > "$COVERAGE_DIR/coverage.out"
for profile in "$PROFILE_DIR"/*.out; do
    tail -n +2 "$profile" >> "$COVERAGE_DIR/coverage.out"
done

# Generar reporte HTML
echo -e "\n${YELLOW}Generando reporte HTML...${NC}"
go tool cover -html="$COVERAGE_DIR/coverage.out" -o "$REPORT_DIR/coverage.html"

# Generar reporte funcional
echo -e "\n${YELLOW}Generando reporte funcional...${NC}"
go tool cover -func="$COVERAGE_DIR/coverage.out" | tee "$REPORT_DIR/coverage.txt"

# Calcular cobertura total
coverage_percent=$(tail -1 "$REPORT_DIR/coverage.txt" | awk '{print substr($3, 1, length($3)-1)}')

# Generar reporte detallado
echo -e "\n${YELLOW}Generando reporte detallado...${NC}"
{
    echo "# Reporte de Cobertura de Código"
    echo "## Resumen"
    echo "- Fecha: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "- Cobertura Total: ${coverage_percent}%"
    echo "- Cobertura Mínima Requerida: ${MIN_COVERAGE}%"
    echo
    echo "## Cobertura por Paquete"
    echo "| Paquete | Cobertura | Estado |"
    echo "|---------|-----------|---------|"
    
    for pkg in $packages; do
        pkg_name=$(basename $pkg)
        pkg_coverage=$(go tool cover -func="$PROFILE_DIR/$pkg_name.out" | tail -1 | awk '{print substr($3, 1, length($3)-1)}')
        if (( $(echo "$pkg_coverage >= $MIN_COVERAGE" | bc -l) )); then
            status="✅"
        else
            status="❌"
        fi
        echo "| $pkg_name | ${pkg_coverage}% | $status |"
    done
} > "$REPORT_DIR/COVERAGE.md"

# Mostrar resumen
echo -e "\n${BLUE}=== Resumen de Cobertura ===${NC}"
echo -e "Cobertura Total: ${YELLOW}${coverage_percent}%${NC}"

if (( $(echo "$coverage_percent >= $MIN_COVERAGE" | bc -l) )); then
    echo -e "${GREEN}✓ La cobertura cumple con el mínimo requerido ($MIN_COVERAGE%)${NC}"
else
    echo -e "${RED}✗ La cobertura está por debajo del mínimo requerido ($MIN_COVERAGE%)${NC}"
fi

# Mostrar ubicación de reportes
echo -e "\n${BLUE}=== Reportes Generados ===${NC}"
echo -e "- Reporte HTML: ${YELLOW}$REPORT_DIR/coverage.html${NC}"
echo -e "- Reporte Funcional: ${YELLOW}$REPORT_DIR/coverage.txt${NC}"
echo -e "- Reporte Markdown: ${YELLOW}$REPORT_DIR/COVERAGE.md${NC}"

# Salir con código de error si no se cumple la cobertura mínima
if (( $(echo "$coverage_percent < $MIN_COVERAGE" | bc -l) )); then
    exit 1
fi

exit 0 