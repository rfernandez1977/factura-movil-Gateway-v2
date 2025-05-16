# Ejemplos de Uso - FMgo

Este documento proporciona ejemplos detallados de uso del módulo FMgo para diferentes casos de uso comunes.

## Índice

1. [Generación de DTEs](#generación-de-dtes)
2. [Validación de Documentos](#validación-de-documentos)
3. [Firma Digital](#firma-digital)
4. [Envío al SII](#envío-al-sii)
5. [Manejo de Respuestas](#manejo-de-respuestas)

## Generación de DTEs

### Factura Electrónica

```go
package main

import (
    "time"
    "github.com/fmgo/core/sii/models"
)

func generarFactura() *models.DTE {
    return &models.DTE{
        Documento: models.Documento{
            Encabezado: models.Encabezado{
                IdDoc: models.IdDoc{
                    TipoDTE: "33",  // 33 = Factura Electrónica
                    Folio:   1,
                    FchEmis: time.Now(),
                },
                Emisor: models.Emisor{
                    RUTEmisor:  "76.123.456-7",
                    RznSoc:     "Empresa Ejemplo S.A.",
                    GiroEmis:   "Servicios Informáticos",
                    Acteco:     "722000",
                    DirOrigen:  "Av. Principal 123",
                    CmnaOrigen: "Santiago",
                },
                Receptor: models.Receptor{
                    RUTRecep:    "77.888.999-0",
                    RznSocRecep: "Cliente Ejemplo Ltda.",
                    GiroRecep:   "Comercio",
                    DirRecep:    "Calle Cliente 456",
                    CmnaRecep:   "Providencia",
                },
                Totales: models.Totales{
                    MntNeto:  100000,
                    TasaIVA:  19.0,
                    IVA:      19000,
                    MntTotal: 119000,
                },
            },
            Detalle: []models.Detalle{
                {
                    NroLinDet: 1,
                    NmbItem:   "Servicio de Desarrollo",
                    QtyItem:   1.0,
                    PrcItem:   100000.0,
                    MontoItem: 100000,
                },
            },
        },
    }
}
```

### Boleta Electrónica

```go
func generarBoleta() *models.DTE {
    return &models.DTE{
        Documento: models.Documento{
            Encabezado: models.Encabezado{
                IdDoc: models.IdDoc{
                    TipoDTE: "39",  // 39 = Boleta Electrónica
                    Folio:   1,
                    FchEmis: time.Now(),
                },
                // ... resto de la configuración ...
            },
        },
    }
}
```

## Validación de Documentos

### Validación Simple

```go
package main

import (
    "log"
    "github.com/fmgo/core/sii/services"
)

func validarDocumento(dte *models.DTE) error {
    validator := services.NewValidatorService()
    if err := validator.CargarEsquemasBase("schema_dte"); err != nil {
        return fmt.Errorf("error cargando esquemas: %w", err)
    }

    return validator.ValidarDTE(dte)
}
```

### Validación con Middleware

```go
func configurarValidacion() {
    middleware := middleware.NewValidatorMiddleware("schema_dte")
    
    // Usar middleware en un handler
    handler := middleware.WithDTEValidation(func(ctx context.Context, dte *models.DTE) error {
        // Procesar DTE validado
        return nil
    })
}
```

### Validación en Lote

```go
func validarLote(dtes []*models.DTE) []error {
    validator := services.NewValidatorService()
    if err := validator.CargarEsquemasBase("schema_dte"); err != nil {
        return []error{err}
    }

    var errores []error
    for i, dte := range dtes {
        if err := validator.ValidarDTE(dte); err != nil {
            errores = append(errores, fmt.Errorf("error en DTE %d: %w", i+1, err))
        }
    }
    return errores
}
```

## Firma Digital

### Firmar DTE Individual

```go
func firmarDocumento(dte *models.DTE) error {
    firmador := services.NewFirmaService(
        os.Getenv("FMGO_CERT_PATH"),
        os.Getenv("FMGO_CERT_PASS"),
        os.Getenv("FMGO_RUT_EMPRESA"),
    )

    return firmador.FirmarDTE(dte)
}
```

### Firmar Lote de DTEs

```go
func firmarLote(dtes []*models.DTE) error {
    firmador := services.NewFirmaService(
        os.Getenv("FMGO_CERT_PATH"),
        os.Getenv("FMGO_CERT_PASS"),
        os.Getenv("FMGO_RUT_EMPRESA"),
    )

    for _, dte := range dtes {
        if err := firmador.FirmarDTE(dte); err != nil {
            return fmt.Errorf("error firmando DTE %s: %w", dte.Documento.IdDoc.Folio, err)
        }
    }
    return nil
}
```

## Envío al SII

### Envío Individual

```go
func enviarAlSII(dte *models.DTE) error {
    client := services.NewSIIClient()
    
    // Autenticar
    if err := client.Autenticar(); err != nil {
        return fmt.Errorf("error de autenticación: %w", err)
    }

    // Enviar
    resp, err := client.EnviarDTE(dte)
    if err != nil {
        return fmt.Errorf("error enviando DTE: %w", err)
    }

    // Procesar respuesta
    return procesarRespuesta(resp)
}
```

### Envío en Lote

```go
func enviarLoteAlSII(dtes []*models.DTE) error {
    client := services.NewSIIClient()
    
    // Crear envío
    envio := services.CrearEnvioDTE(dtes)
    
    // Enviar
    resp, err := client.EnviarLote(envio)
    if err != nil {
        return fmt.Errorf("error enviando lote: %w", err)
    }

    return procesarRespuestaLote(resp)
}
```

## Manejo de Respuestas

### Procesar Respuesta

```go
func procesarRespuesta(resp *models.RespuestaSII) error {
    if !resp.EsExitosa() {
        return fmt.Errorf("error en envío: %s", resp.Mensaje)
    }

    // Guardar trackID
    trackID := resp.TrackID
    
    // Consultar estado
    estado, err := consultarEstado(trackID)
    if err != nil {
        return fmt.Errorf("error consultando estado: %w", err)
    }

    return nil
}
```

### Consultar Estado

```go
func consultarEstado(trackID string) (*models.EstadoEnvio, error) {
    client := services.NewSIIClient()
    
    estado, err := client.ConsultarEstado(trackID)
    if err != nil {
        return nil, fmt.Errorf("error consultando estado: %w", err)
    }

    return estado, nil
}
```

## Ejemplos Completos

### Flujo Completo de Facturación

```go
func emitirFactura() error {
    // 1. Generar DTE
    dte := generarFactura()

    // 2. Validar
    validator := services.NewValidatorService()
    if err := validator.CargarEsquemasBase("schema_dte"); err != nil {
        return err
    }
    if err := validator.ValidarDTE(dte); err != nil {
        return err
    }

    // 3. Firmar
    firmador := services.NewFirmaService(
        os.Getenv("FMGO_CERT_PATH"),
        os.Getenv("FMGO_CERT_PASS"),
        os.Getenv("FMGO_RUT_EMPRESA"),
    )
    if err := firmador.FirmarDTE(dte); err != nil {
        return err
    }

    // 4. Enviar al SII
    client := services.NewSIIClient()
    resp, err := client.EnviarDTE(dte)
    if err != nil {
        return err
    }

    // 5. Procesar respuesta
    return procesarRespuesta(resp)
}
```

### Manejo de Errores

```go
func manejarErrores(err error) {
    switch e := err.(type) {
    case *services.ValidacionError:
        log.Printf("Error de validación: %v", e)
        // Manejar error de validación
    case *services.FirmaError:
        log.Printf("Error de firma: %v", e)
        // Manejar error de firma
    case *services.EnvioError:
        log.Printf("Error de envío: %v", e)
        // Manejar error de envío
    default:
        log.Printf("Error general: %v", err)
        // Manejar otros errores
    }
}
```

## Buenas Prácticas

1. **Validación Temprana**: Siempre validar los documentos antes de firmarlos y enviarlos.
2. **Manejo de Errores**: Implementar un manejo de errores robusto y específico.
3. **Logging**: Mantener un registro detallado de todas las operaciones.
4. **Reintentos**: Implementar lógica de reintentos para operaciones de red.
5. **Caché**: Utilizar el caché de tokens para optimizar las autenticaciones.

## Referencias

- [Documentación API SII](https://www.sii.cl/factura_electronica/factura_mercado/api_ref.html)
- [Formatos XML](https://www.sii.cl/factura_electronica/factura_mercado/formato_dte.pdf)
- [Códigos de Error](https://www.sii.cl/factura_electronica/factura_mercado/codigos_error.html) 