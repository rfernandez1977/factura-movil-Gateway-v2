#!/bin/bash

# Script de configuración del ambiente de certificación
# Este script configura el ambiente necesario para las pruebas de certificación con el SII

# Colores para mensajes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Función para imprimir mensajes de estado
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Verificar que estamos en el directorio correcto
if [ ! -d "dev" ] || [ ! -d "test" ]; then
    print_error "Este script debe ejecutarse desde el directorio raíz del proyecto"
    exit 1
fi

# Crear estructura de directorios si no existe
mkdir -p dev/config/{caf,certs,temp,logs}
mkdir -p test/config/{caf,certs,temp,logs}
print_status "Estructura de directorios creada"

# Establecer permisos
chmod 755 dev/config
chmod 700 dev/config/certs
chmod 700 dev/config/caf
chmod 755 dev/config/temp
chmod 755 dev/config/logs
chmod 755 test/config
chmod 700 test/config/certs
chmod 700 test/config/caf
chmod 755 test/config/temp
chmod 755 test/config/logs
print_status "Permisos establecidos"

# Verificar certificados
if [ ! -f "dev/config/certs/cert.p12" ]; then
    print_error "Certificado de desarrollo no encontrado en dev/config/certs/cert.p12"
    echo "Por favor, copie su certificado de desarrollo a dev/config/certs/cert.p12"
fi

if [ ! -f "test/config/certs/cert.p12" ]; then
    print_error "Certificado de pruebas no encontrado en test/config/certs/cert.p12"
    echo "Por favor, copie su certificado de pruebas a test/config/certs/cert.p12"
fi

# Verificar archivos CAF
if [ ! -f "dev/config/caf/FoliosSII.xml" ]; then
    print_error "Archivo CAF de desarrollo no encontrado en dev/config/caf/FoliosSII.xml"
    echo "Por favor, copie su archivo CAF de desarrollo a dev/config/caf/FoliosSII.xml"
fi

if [ ! -f "test/config/caf/FoliosSII.xml" ]; then
    print_error "Archivo CAF de pruebas no encontrado en test/config/caf/FoliosSII.xml"
    echo "Por favor, copie su archivo CAF de pruebas a test/config/caf/FoliosSII.xml"
fi

# Crear archivo de configuración si no existe
cat > dev/config/config.json << EOF
{
    "certificado": {
        "ruta": "dev/config/certs/cert.p12",
        "password": "",
        "ttl_cache": 3600,
        "max_items_cache": 100
    },
    "caf": {
        "ruta": "dev/config/caf/FoliosSII.xml"
    },
    "sii": {
        "ambiente": "certificacion",
        "url_base": "https://maullin.sii.cl/DTEWS/",
        "url_semilla": "CrSeed.jws?WSDL",
        "url_token": "GetTokenFromSeed.jws?WSDL"
    },
    "logs": {
        "nivel": "debug",
        "ruta": "dev/config/logs/app.log"
    }
}
EOF
print_status "Archivo de configuración creado en dev/config/config.json"

# Copiar configuración a ambiente de pruebas
cp dev/config/config.json test/config/config.json
sed -i '' 's/dev\/config/test\/config/g' test/config/config.json
print_status "Archivo de configuración copiado a test/config/config.json"

print_status "Configuración del ambiente completada"
echo "Por favor, complete los siguientes pasos manualmente:"
echo "1. Copie su certificado digital a dev/config/certs/cert.p12"
echo "2. Copie su archivo CAF a dev/config/caf/FoliosSII.xml"
echo "3. Configure la contraseña del certificado en dev/config/config.json"
echo "4. Repita los pasos 1-3 para el ambiente de pruebas (test/config/)" 