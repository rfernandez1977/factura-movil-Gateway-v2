#!/bin/bash

# Buscar todos los archivos .go y actualizar las importaciones
find . -type f -name "*.go" -exec sed -i '' 's|"github.com/fmgo/|"FMgo/|g' {} +

# Actualizar el User-Agent en sesion_service.go
sed -i '' 's|github.com/fmgo/1.0|FMgo/1.0|g' services/sesion_service.go

# Mostrar los archivos modificados
echo "Archivos actualizados:"
git diff --name-only

# Actualizar go.mod
go mod tidy 