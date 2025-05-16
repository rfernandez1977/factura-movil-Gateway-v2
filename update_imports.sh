#!/bin/bash

# Actualizar importaciones en todo el proyecto
find . -type f -name "*.go" -exec sed -i '' 's|github.com/cursor/FMgo|github.com/fmgo|g' {} +

# Mostrar los archivos modificados
echo "Archivos actualizados:"
git diff --name-only

# Actualizar go.mod
go mod tidy 