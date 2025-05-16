# Validador CAF - MVP

## Descripción
El validador CAF (Código de Autorización de Folios) es un componente esencial para la gestión de folios en documentos tributarios electrónicos. Esta implementación MVP proporciona las funcionalidades básicas necesarias para validar y controlar el uso de folios.

## Características Implementadas

### Validaciones Básicas
- Validación de RUT emisor
- Validación de tipo de DTE
- Control de rango de folios
- Validación de fechas de vigencia
- Control de folios usados (en memoria)

### Estructura de Datos
```go
type CAF struct {
    RutEmisor  string    
    TipoDTE    int       
    RangoDesde int       
    RangoHasta int       
    FechaDesde time.Time 
    FechaHasta time.Time 
}
```

### Manejo de Errores
- `ErrCAFInvalido`: CAF mal formado o inválido
- `ErrCAFExpirado`: CAF fuera de vigencia
- `ErrFolioNoValido`: Folio fuera de rango
- `ErrRUTNoCoincide`: RUT no coincide
- `ErrTipoDTEInvalido`: Tipo de DTE incorrecto
- `ErrFolioUsado`: Folio ya utilizado

## Uso

### Crear Validador
```go
validator, err := caf.NewValidator(cafXMLData)
if err != nil {
    log.Fatal(err)
}
```

### Validar Folio
```go
// Validación individual
if err := validator.ValidarFolio(123); err != nil {
    log.Printf("Folio inválido: %v", err)
}

// Validación completa
if err := validator.ValidarCompleto("76212889-6", 33, 123); err != nil {
    log.Printf("Error en validación: %v", err)
}
```

### Control de Folios
```go
// Marcar folio como usado
if err := validator.MarcarFolioUsado(123); err != nil {
    log.Printf("Error marcando folio: %v", err)
}
```

## Formato XML CAF
```xml
<?xml version="1.0" encoding="UTF-8"?>
<AUTORIZACION>
    <CAF>
        <DA>
            <RE>76212889-6</RE>
            <TD>33</TD>
            <RNG>
                <D>1</D>
                <H>100</H>
            </RNG>
            <RSAPK>
                <M>2024-03-01T00:00:00Z</M>
                <E>2024-12-31T23:59:59Z</E>
            </RSAPK>
        </DA>
    </CAF>
</AUTORIZACION>
```

## Características Post-MVP

### Verificación de Firmas
- Implementación de verificación de firmas XML
- Validación de certificados
- Manejo de claves públicas

### Persistencia
- Almacenamiento persistente de folios usados
- Sincronización entre instancias
- Recuperación de estado

### Monitoreo
- Métricas de uso
- Alertas de agotamiento
- Logs estructurados

### Optimizaciones
- Caché de validaciones
- Manejo de concurrencia
- Pruebas de carga

## Pruebas
Las pruebas unitarias cubren los siguientes escenarios:
- Validación de CAF válido/inválido
- Validación de rangos de folios
- Validación de RUT emisor
- Validación de tipo DTE
- Control de folios usados

## Limitaciones Actuales
- Sin persistencia de folios usados
- Sin verificación de firmas
- Sin manejo de concurrencia avanzado
- Sin métricas ni monitoreo

## Referencias
- [Documentación SII - CAF](https://www.sii.cl/factura_electronica/factura_mercado/CAF.pdf)
- [Especificación XML](https://www.sii.cl/factura_electronica/factura_mercado/XML.pdf) 