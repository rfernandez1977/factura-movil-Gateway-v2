# Validador CAF

Este paquete implementa la validación de Códigos de Autorización de Folios (CAF) para el sistema de facturación electrónica.

## Características MVP

### Validaciones Básicas
- Validación de RUT emisor
- Validación de tipo DTE
- Control de rango de folios
- Validación de fechas de vigencia
- Control de folios usados (en memoria)

### Servicio de Gestión
- Registro de CAFs
- Validación de folios
- Consulta de estado

## Uso

### Registrar un CAF
```go
service := caf.NewService()
err := service.RegistrarCAF(cafXMLData)
if err != nil {
    log.Fatal(err)
}
```

### Validar un Folio
```go
err := service.ValidarFolio("76212889-6", 33, 50)
if err != nil {
    log.Printf("Error validando folio: %v", err)
}
```

### Consultar Estado
```go
estado, err := service.ObtenerEstadoCAF("76212889-6", 33)
if err != nil {
    log.Printf("Error obteniendo estado: %v", err)
}
fmt.Printf("Rango de folios: %d-%d\n", estado.RangoDesde, estado.RangoHasta)
```

## Estructura XML CAF
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
                <M>2023-01-01T00:00:00Z</M>
                <E>2025-12-31T23:59:59Z</E>
            </RSAPK>
        </DA>
    </CAF>
</AUTORIZACION>
```

## Errores
- `ErrCAFInvalido`: CAF mal formado o inválido
- `ErrCAFExpirado`: CAF fuera de vigencia
- `ErrFolioNoValido`: Folio fuera de rango
- `ErrRUTNoCoincide`: RUT no coincide
- `ErrTipoDTEInvalido`: Tipo de DTE incorrecto
- `ErrFolioUsado`: Folio ya utilizado

## Próximas Características
- Verificación de firmas XML
- Persistencia de folios usados
- Métricas y monitoreo
- Pruebas de carga
- Manejo avanzado de concurrencia

## Pruebas
```bash
go test -v ./...
```

## Contribución
1. Fork el proyecto
2. Crear rama feature (`git checkout -b feature/NuevaFuncionalidad`)
3. Commit cambios (`git commit -m 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/NuevaFuncionalidad`)
5. Crear Pull Request 