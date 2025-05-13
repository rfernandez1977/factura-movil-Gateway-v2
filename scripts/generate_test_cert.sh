#!/bin/bash

# Crear directorio para el certificado
mkdir -p testdata

# Generar llave privada
openssl genrsa -out testdata/private.key 2048

# Generar certificado autofirmado
openssl req -new -x509 -key testdata/private.key -out testdata/certificate.crt -days 365 \
    -subj "/C=CL/ST=Metropolitana/L=Santiago/O=Empresa Test/CN=76.123.456-7"

# Convertir certificado y llave a formato PFX
openssl pkcs12 -export -out testdata/certificado.pfx \
    -inkey testdata/private.key \
    -in testdata/certificate.crt \
    -passout pass:contraseña

# Limpiar archivos intermedios
rm testdata/private.key testdata/certificate.crt 