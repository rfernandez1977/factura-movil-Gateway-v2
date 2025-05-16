#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Instalando herramientas de desarrollo...${NC}"

# Verificar Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go no está instalado. Por favor, instale Go 1.24.2 o superior.${NC}"
    exit 1
fi

# Instalar herramientas de desarrollo
echo -e "${YELLOW}Instalando golangci-lint...${NC}"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2

echo -e "${YELLOW}Instalando gotestsum...${NC}"
go install gotest.tools/gotestsum@v1.11.0

echo -e "${YELLOW}Instalando mockgen...${NC}"
go install github.com/golang/mock/mockgen@v1.6.0

echo -e "${YELLOW}Instalando swag...${NC}"
go install github.com/swaggo/swag/cmd/swag@v1.16.3

echo -e "${YELLOW}Instalando staticcheck...${NC}"
go install honnef.co/go/tools/cmd/staticcheck@latest

echo -e "${YELLOW}Instalando gosec...${NC}"
go install github.com/securego/gosec/v2/cmd/gosec@latest

echo -e "${YELLOW}Instalando nancy...${NC}"
go install github.com/sonatype-nexus-community/nancy@latest

echo -e "${YELLOW}Instalando migrate...${NC}"
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

echo -e "${YELLOW}Instalando goose...${NC}"
go install github.com/pressly/goose/v3/cmd/goose@latest

# Verificar instalaciones
echo -e "${YELLOW}Verificando instalaciones...${NC}"

tools=(
    "golangci-lint"
    "gotestsum"
    "mockgen"
    "swag"
    "staticcheck"
    "gosec"
    "nancy"
    "migrate"
    "goose"
)

for tool in "${tools[@]}"; do
    if command -v $tool &> /dev/null; then
        echo -e "${GREEN}✓ $tool instalado correctamente${NC}"
    else
        echo -e "${RED}✗ Error al instalar $tool${NC}"
    fi
done

echo -e "${GREEN}¡Instalación completada!${NC}" 