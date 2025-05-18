#!/bin/bash

# Script para generar folios CAF de prueba
# Este script genera archivos CAF para diferentes tipos de documentos

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
mkdir -p test/data/caf/{active,expired,depleted}
print_status "Directorios CAF creados"

# Generar CAF activo para factura (tipo 33)
cat > test/data/caf/active/CAF_33.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF version="1.0">
        <DA>
            <RE>76555555-5</RE>
            <RS>EMPRESA DE PRUEBA SPA</RS>
            <TD>33</TD>
            <RNG>
                <D>1</D>
                <H>100</H>
            </RNG>
            <FA>2024-03-23</FA>
            <RSAPK>
                <M>1234567890</M>
                <E>010001</E>
            </RSAPK>
            <IDK>100</IDK>
        </DA>
        <FRMA>ABC123</FRMA>
    </CAF>
    <RSASK>DEF456</RSASK>
    <RSAPUBK>GHI789</RSAPUBK>
</AUTORIZACION>
EOF
print_status "CAF activo para factura generado"

# Generar CAF activo para nota de crédito (tipo 61)
cat > test/data/caf/active/CAF_61.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF version="1.0">
        <DA>
            <RE>76555555-5</RE>
            <RS>EMPRESA DE PRUEBA SPA</RS>
            <TD>61</TD>
            <RNG>
                <D>1</D>
                <H>100</H>
            </RNG>
            <FA>2024-03-23</FA>
            <RSAPK>
                <M>1234567890</M>
                <E>010001</E>
            </RSAPK>
            <IDK>101</IDK>
        </DA>
        <FRMA>ABC123</FRMA>
    </CAF>
    <RSASK>DEF456</RSASK>
    <RSAPUBK>GHI789</RSAPUBK>
</AUTORIZACION>
EOF
print_status "CAF activo para nota de crédito generado"

# Generar CAF activo para boleta (tipo 39)
cat > test/data/caf/active/CAF_39.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF version="1.0">
        <DA>
            <RE>76555555-5</RE>
            <RS>EMPRESA DE PRUEBA SPA</RS>
            <TD>39</TD>
            <RNG>
                <D>1</D>
                <H>100</H>
            </RNG>
            <FA>2024-03-23</FA>
            <RSAPK>
                <M>1234567890</M>
                <E>010001</E>
            </RSAPK>
            <IDK>102</IDK>
        </DA>
        <FRMA>ABC123</FRMA>
    </CAF>
    <RSASK>DEF456</RSASK>
    <RSAPUBK>GHI789</RSAPUBK>
</AUTORIZACION>
EOF
print_status "CAF activo para boleta generado"

# Generar CAF expirado
cat > test/data/caf/expired/CAF_33.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF version="1.0">
        <DA>
            <RE>76555555-5</RE>
            <RS>EMPRESA DE PRUEBA SPA</RS>
            <TD>33</TD>
            <RNG>
                <D>1</D>
                <H>100</H>
            </RNG>
            <FA>2023-03-23</FA>
            <RSAPK>
                <M>1234567890</M>
                <E>010001</E>
            </RSAPK>
            <IDK>103</IDK>
        </DA>
        <FRMA>ABC123</FRMA>
    </CAF>
    <RSASK>DEF456</RSASK>
    <RSAPUBK>GHI789</RSAPUBK>
</AUTORIZACION>
EOF
print_status "CAF expirado generado"

# Generar CAF agotado
cat > test/data/caf/depleted/CAF_33.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF version="1.0">
        <DA>
            <RE>76555555-5</RE>
            <RS>EMPRESA DE PRUEBA SPA</RS>
            <TD>33</TD>
            <RNG>
                <D>1</D>
                <H>1</H>
            </RNG>
            <FA>2024-03-23</FA>
            <RSAPK>
                <M>1234567890</M>
                <E>010001</E>
            </RSAPK>
            <IDK>104</IDK>
        </DA>
        <FRMA>ABC123</FRMA>
    </CAF>
    <RSASK>DEF456</RSASK>
    <RSAPUBK>GHI789</RSAPUBK>
</AUTORIZACION>
EOF
print_status "CAF agotado generado"

# Copiar CAF activo a los directorios de configuración
cp test/data/caf/active/CAF_33.xml dev/config/caf/FoliosSII.xml
cp test/data/caf/active/CAF_33.xml test/config/caf/FoliosSII.xml
print_status "CAF activo copiado a directorios de configuración"

print_status "Generación de CAF completada"
echo "Los archivos CAF se han generado en test/data/caf/"
echo "Por favor, verifique los siguientes directorios:"
echo "- test/data/caf/active/"
echo "- test/data/caf/expired/"
echo "- test/data/caf/depleted/" 