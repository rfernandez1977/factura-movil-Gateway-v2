#!/bin/bash

# Script para generar datos de prueba para certificación
# Este script genera documentos XML de prueba para diferentes tipos de DTE

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

# Crear directorio para datos de prueba
mkdir -p test/data/dte
print_status "Directorio de datos de prueba creado"

# Generar factura electrónica de prueba
cat > test/data/dte/factura.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <Documento ID="F1">
        <Encabezado>
            <IdDoc>
                <TipoDTE>33</TipoDTE>
                <Folio>1</Folio>
                <FchEmis>2024-03-23</FchEmis>
            </IdDoc>
            <Emisor>
                <RUTEmisor>76555555-5</RUTEmisor>
                <RznSoc>EMPRESA DE PRUEBA SPA</RznSoc>
                <GiroEmis>DESARROLLO DE SOFTWARE</GiroEmis>
                <Acteco>722000</Acteco>
                <DirOrigen>CALLE PRUEBA 123</DirOrigen>
                <CmnaOrigen>SANTIAGO</CmnaOrigen>
            </Emisor>
            <Receptor>
                <RUTRecep>55666777-8</RUTRecep>
                <RznSocRecep>CLIENTE DE PRUEBA LTDA</RznSocRecep>
                <GiroRecep>COMERCIO</GiroRecep>
                <DirRecep>AVENIDA TEST 456</DirRecep>
                <CmnaRecep>PROVIDENCIA</CmnaRecep>
            </Receptor>
            <Totales>
                <MntNeto>100000</MntNeto>
                <TasaIVA>19</TasaIVA>
                <IVA>19000</IVA>
                <MntTotal>119000</MntTotal>
            </Totales>
        </Encabezado>
        <Detalle>
            <NroLinDet>1</NroLinDet>
            <NmbItem>Servicio de Desarrollo</NmbItem>
            <QtyItem>1</QtyItem>
            <PrcItem>100000</PrcItem>
            <MontoItem>100000</MontoItem>
        </Detalle>
    </Documento>
</DTE>
EOF
print_status "Factura de prueba generada"

# Generar nota de crédito de prueba
cat > test/data/dte/nota_credito.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <Documento ID="NC1">
        <Encabezado>
            <IdDoc>
                <TipoDTE>61</TipoDTE>
                <Folio>1</Folio>
                <FchEmis>2024-03-23</FchEmis>
            </IdDoc>
            <Emisor>
                <RUTEmisor>76555555-5</RUTEmisor>
                <RznSoc>EMPRESA DE PRUEBA SPA</RznSoc>
                <GiroEmis>DESARROLLO DE SOFTWARE</GiroEmis>
                <Acteco>722000</Acteco>
                <DirOrigen>CALLE PRUEBA 123</DirOrigen>
                <CmnaOrigen>SANTIAGO</CmnaOrigen>
            </Emisor>
            <Receptor>
                <RUTRecep>55666777-8</RUTRecep>
                <RznSocRecep>CLIENTE DE PRUEBA LTDA</RznSocRecep>
                <GiroRecep>COMERCIO</GiroRecep>
                <DirRecep>AVENIDA TEST 456</DirRecep>
                <CmnaRecep>PROVIDENCIA</CmnaRecep>
            </Receptor>
            <Totales>
                <MntNeto>-50000</MntNeto>
                <TasaIVA>19</TasaIVA>
                <IVA>-9500</IVA>
                <MntTotal>-59500</MntTotal>
            </Totales>
        </Encabezado>
        <Detalle>
            <NroLinDet>1</NroLinDet>
            <NmbItem>Descuento Servicio</NmbItem>
            <QtyItem>1</QtyItem>
            <PrcItem>50000</PrcItem>
            <MontoItem>50000</MontoItem>
        </Detalle>
        <Referencia>
            <NroLinRef>1</NroLinRef>
            <TpoDocRef>33</TpoDocRef>
            <FolioRef>1</FolioRef>
            <FchRef>2024-03-23</FchRef>
            <CodRef>1</CodRef>
            <RazonRef>DESCUENTO POR PRONTO PAGO</RazonRef>
        </Referencia>
    </Documento>
</DTE>
EOF
print_status "Nota de crédito de prueba generada"

# Generar boleta electrónica de prueba
cat > test/data/dte/boleta.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <Documento ID="B1">
        <Encabezado>
            <IdDoc>
                <TipoDTE>39</TipoDTE>
                <Folio>1</Folio>
                <FchEmis>2024-03-23</FchEmis>
            </IdDoc>
            <Emisor>
                <RUTEmisor>76555555-5</RUTEmisor>
                <RznSoc>EMPRESA DE PRUEBA SPA</RznSoc>
                <GiroEmis>DESARROLLO DE SOFTWARE</GiroEmis>
                <Acteco>722000</Acteco>
                <DirOrigen>CALLE PRUEBA 123</DirOrigen>
                <CmnaOrigen>SANTIAGO</CmnaOrigen>
            </Emisor>
            <Receptor>
                <RUTRecep>66777888-9</RUTRecep>
            </Receptor>
            <Totales>
                <MntTotal>25000</MntTotal>
            </Totales>
        </Encabezado>
        <Detalle>
            <NroLinDet>1</NroLinDet>
            <NmbItem>Soporte Técnico</NmbItem>
            <QtyItem>1</QtyItem>
            <PrcItem>25000</PrcItem>
            <MontoItem>25000</MontoItem>
        </Detalle>
    </Documento>
</DTE>
EOF
print_status "Boleta de prueba generada"

# Generar datos de prueba para semilla
cat > test/data/dte/semilla.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<SemillaXML>
    <Semilla>1234567890</Semilla>
</SemillaXML>
EOF
print_status "Datos de semilla generados"

# Generar datos de prueba para token
cat > test/data/dte/token.xml << EOF
<?xml version="1.0" encoding="UTF-8"?>
<TokenXML>
    <Token>ABC123DEF456GHI789</Token>
</TokenXML>
EOF
print_status "Datos de token generados"

print_status "Generación de datos de prueba completada"
echo "Los archivos de prueba se han generado en test/data/dte/"
echo "Por favor, verifique los siguientes archivos:"
echo "- test/data/dte/factura.xml"
echo "- test/data/dte/nota_credito.xml"
echo "- test/data/dte/boleta.xml"
echo "- test/data/dte/semilla.xml"
echo "- test/data/dte/token.xml" 