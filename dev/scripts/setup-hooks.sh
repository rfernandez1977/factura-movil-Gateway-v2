#!/bin/bash

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

HOOKS_DIR=".git/hooks"
CUSTOM_HOOKS_DIR="dev/scripts/git-hooks"

echo -e "${YELLOW}Configurando git hooks...${NC}"

# Crear directorio de hooks personalizados si no existe
mkdir -p "$CUSTOM_HOOKS_DIR"

# Crear pre-commit hook
cat > "$CUSTOM_HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash

echo "Ejecutando pre-commit hooks..."

# Ejecutar gofmt
echo "Verificando formato del código..."
if ! gofmt -l . | grep -q '^'; then
    echo "✓ Formato del código correcto"
else
    echo "✗ Error: Hay archivos que necesitan ser formateados"
    gofmt -l .
    exit 1
fi

# Ejecutar golangci-lint
echo "Ejecutando linter..."
if golangci-lint run; then
    echo "✓ Lint pasó correctamente"
else
    echo "✗ Error: El linter encontró problemas"
    exit 1
fi

# Ejecutar tests
echo "Ejecutando tests..."
if go test ./...; then
    echo "✓ Tests pasaron correctamente"
else
    echo "✗ Error: Algunos tests fallaron"
    exit 1
fi

# Ejecutar gosec
echo "Ejecutando análisis de seguridad..."
if gosec ./...; then
    echo "✓ Análisis de seguridad pasó correctamente"
else
    echo "✗ Error: Se encontraron problemas de seguridad"
    exit 1
fi

echo "✓ Todos los checks pre-commit pasaron correctamente"
exit 0
EOF

# Crear pre-push hook
cat > "$CUSTOM_HOOKS_DIR/pre-push" << 'EOF'
#!/bin/bash

echo "Ejecutando pre-push hooks..."

# Ejecutar tests completos
echo "Ejecutando suite completa de tests..."
if go test ./...; then
    echo "✓ Tests pasaron correctamente"
else
    echo "✗ Error: Algunos tests fallaron"
    exit 1
fi

# Ejecutar tests de race conditions
echo "Ejecutando tests de race conditions..."
if go test -race ./...; then
    echo "✓ No se detectaron race conditions"
else
    echo "✗ Error: Se detectaron race conditions"
    exit 1
fi

echo "✓ Todos los checks pre-push pasaron correctamente"
exit 0
EOF

# Hacer ejecutables los hooks
chmod +x "$CUSTOM_HOOKS_DIR/pre-commit"
chmod +x "$CUSTOM_HOOKS_DIR/pre-push"

# Crear enlaces simbólicos
ln -sf "../../$CUSTOM_HOOKS_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
ln -sf "../../$CUSTOM_HOOKS_DIR/pre-push" "$HOOKS_DIR/pre-push"

echo -e "${GREEN}¡Git hooks configurados correctamente!${NC}"
echo -e "${YELLOW}Hooks instalados:${NC}"
echo "- pre-commit"
echo "- pre-push" 