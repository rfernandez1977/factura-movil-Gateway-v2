#!/bin/bash

# Actualizar todas las importaciones que usan FMgo a github.com/fmgo
find . -type f -name "*.go" -exec sed -i '' 's/"FMgo\//"github.com\/fmgo\//g' {} +

# Actualizar las importaciones en los archivos de prueba
find ./tests -type f -name "*_test.go" -exec sed -i '' 's/package main/package integration/g' {} +

# Ejecutar go mod tidy para limpiar y actualizar dependencias
go mod tidy 