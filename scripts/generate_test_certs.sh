#!/bin/bash

# Script para generar certificados de prueba
# Este script genera certificados para pruebas de firma digital

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
if [ ! -d "test/config" ]; then
    print_error "Este script debe ejecutarse desde el directorio raíz del proyecto"
    exit 1
fi

# Crear directorios necesarios
mkdir -p test/data/certs/{valid,expired,revoked}
print_status "Directorios de certificados creados"

# Generar certificado válido
openssl req -x509 -newkey rsa:2048 -keyout test/data/certs/valid/key.pem -out test/data/certs/valid/cert.pem -days 365 -nodes -subj "/C=CL/ST=RM/L=Santiago/O=Empresa de Prueba/OU=Desarrollo/CN=valid.test.local"
print_status "Certificado válido generado"

# Generar certificado expirado (fecha en el pasado)
openssl req -x509 -newkey rsa:2048 -keyout test/data/certs/expired/key.pem -out test/data/certs/expired/cert.pem -days -365 -nodes -subj "/C=CL/ST=RM/L=Santiago/O=Empresa de Prueba/OU=Desarrollo/CN=expired.test.local"
print_status "Certificado expirado generado"

# Generar certificado para revocar
openssl req -x509 -newkey rsa:2048 -keyout test/data/certs/revoked/key.pem -out test/data/certs/revoked/cert.pem -days 365 -nodes -subj "/C=CL/ST=RM/L=Santiago/O=Empresa de Prueba/OU=Desarrollo/CN=revoked.test.local"
print_status "Certificado para revocar generado"

# Convertir certificados a formato P12
openssl pkcs12 -export -out test/data/certs/valid/cert.p12 -inkey test/data/certs/valid/key.pem -in test/data/certs/valid/cert.pem -passout pass:test123
openssl pkcs12 -export -out test/data/certs/expired/cert.p12 -inkey test/data/certs/expired/key.pem -in test/data/certs/expired/cert.pem -passout pass:test123
openssl pkcs12 -export -out test/data/certs/revoked/cert.p12 -inkey test/data/certs/revoked/key.pem -in test/data/certs/revoked/cert.pem -passout pass:test123
print_status "Certificados convertidos a formato P12"

# Copiar certificado válido a los directorios de configuración
cp test/data/certs/valid/cert.p12 dev/config/certs/cert.p12
cp test/data/certs/valid/cert.p12 test/config/certs/cert.p12
print_status "Certificado válido copiado a directorios de configuración"

print_status "Generación de certificados completada"
echo "Los certificados se han generado en test/data/certs/"
echo "Contraseña para todos los certificados P12: test123"
echo "Por favor, verifique los siguientes directorios:"
echo "- test/data/certs/valid/"
echo "- test/data/certs/expired/"
echo "- test/data/certs/revoked/" 