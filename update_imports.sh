#!/bin/bash

# Buscar y reemplazar en todos los archivos .go
find . -type f -name "*.go" -exec sed -i '' 's|github.com/rodrigofernandezcalderon/FMgo|github.com/cursor/FMgo|g' {} +

# Actualizar go.mod
go mod tidy 