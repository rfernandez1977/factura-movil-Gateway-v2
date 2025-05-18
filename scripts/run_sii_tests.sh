#!/bin/bash

echo "🔄 Ejecutando pruebas de conexión con el SII..."

# Crear directorios necesarios
echo "📁 Creando directorios necesarios..."
mkdir -p bin config/certs

# Verificar archivo de configuración
if [ ! -f config/sii_config.json ]; then
    echo "❌ Error: No se encuentra el archivo de configuración config/sii_config.json"
    exit 1
fi

# Compilar el script de prueba
echo "📦 Compilando script de prueba..."
go build -o bin/test_sii_connection scripts/test_sii_connection.go

if [ $? -ne 0 ]; then
    echo "❌ Error compilando el script"
    exit 1
fi

# Ejecutar las pruebas
echo -e "\n🚀 Ejecutando pruebas de conexión..."
./bin/test_sii_connection

# Verificar resultado
if [ $? -eq 0 ]; then
    echo -e "\n✅ Pruebas de conexión completadas exitosamente"
else
    echo -e "\n❌ Error en las pruebas de conexión"
    exit 1
fi 