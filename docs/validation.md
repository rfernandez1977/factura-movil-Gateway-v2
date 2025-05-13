# Sistema de Validación de DTEs

Este documento describe el sistema de validación implementado para documentos tributarios electrónicos (DTEs).

## Índice
1. [Introducción](#introducción)
2. [Componentes Principales](#componentes-principales)
3. [Ejemplos de Uso](#ejemplos-de-uso)
4. [Mejores Prácticas](#mejores-prácticas)

## Introducción

El sistema de validación proporciona una estructura robusta para validar documentos tributarios electrónicos, asegurando que cumplan con los requisitos del SII (Servicio de Impuestos Internos de Chile).

## Componentes Principales

### ValidationRule

```go
rule := NewValidationRule(
    "ValidaRUT",
    "Valida el formato y dígito verificador del RUT",
    "RUT",
    "^[0-9]+-[0-9kK]$",
    "RUT inválido"
)
```

### ValidationConfig

```go
config := NewValidationConfig(
    "DTE_FACTURA",
    []ValidationRule{rule},
    5,           // máximo de errores
    false        // no detener en primer error
)
```

### Ejemplo de Validación

```go
// Crear una solicitud de validación
request := ValidationRequest{
    DocumentoID: "DOC001",
    Tipo:       "FACTURA",
    Config:     *config,
    Metadata: map[string]interface{}{
        "version": "1.0",
        "emisor":  "76.123.456-7",
    },
}

// Procesar la validación
response, err := validator.Validate(request)
if err != nil {
    log.Printf("Error en validación: %v", err)
    return
}

// Verificar resultados
if !response.Exitoso {
    for _, err := range response.Errores {
        log.Printf("Error en campo %s: %s", err.Field, err.Message)
    }
}
```

## Mejores Prácticas

### 1. Definición de Reglas

- Usar nombres descriptivos
- Incluir mensajes claros
- Mantener expresiones simples
- Documentar el propósito

### 2. Configuración

```go
// Ejemplo de configuración recomendada
config := ValidationConfig{
    Tipo:        "FACTURA",
    MaxErrores:  10,
    StopOnError: true,
    Reglas: []ValidationRule{
        {
            Nombre:    "ValidaRUT",
            Tipo:      "RUT",
            Mensaje:   "RUT inválido",
            Activo:    true,
        },
        {
            Nombre:    "ValidaEmail",
            Tipo:      "EMAIL",
            Mensaje:   "Email inválido",
            Activo:    true,
        },
    },
}
```

### 3. Manejo de Errores

```go
// Ejemplo de manejo de errores recomendado
if err := documento.Validate(); err != nil {
    if validErr, ok := err.(*ValidationError); ok {
        log.Printf("Error de validación en campo %s: %s",
            validErr.Field,
            validErr.Message)
        return
    }
    log.Printf("Error inesperado: %v", err)
    return
}
```

## Casos de Uso Comunes

### 1. Validación de Factura

```go
func ValidateFactura(factura *Factura) error {
    validator := NewBaseValidator()
    
    // Validar RUT emisor
    if err := ValidateRUT(factura.Emisor.RUT); err != nil {
        validator.AddError("emisor_rut", err.Error(), "INVALID_RUT")
    }
    
    // Validar correo
    if err := ValidateEmail(factura.Emisor.Email); err != nil {
        validator.AddError("emisor_email", err.Error(), "INVALID_EMAIL")
    }
    
    if validator.HasErrors() {
        return validator.GetErrors()[0]
    }
    return nil
}
```

### 2. Validación de Boleta

```go
func ValidateBoleta(boleta *Boleta) error {
    config := NewValidationConfig(
        "BOLETA",
        []ValidationRule{
            *NewValidationRule("ValidaTotal", "Valida monto total", "MONTO", "", "Monto inválido"),
        },
        3,
        false,
    )
    
    request := ValidationRequest{
        DocumentoID: boleta.ID,
        Tipo:       "BOLETA",
        Config:     *config,
    }
    
    return ProcessValidation(request)
}
```

## Notas Importantes

1. Siempre verificar la versión del documento
2. Mantener reglas actualizadas según normativa SII
3. Documentar cambios en reglas de validación
4. Realizar pruebas exhaustivas

## Contribución

Para contribuir nuevas reglas de validación:

1. Crear la regla en formato estándar
2. Documentar el propósito y uso
3. Incluir pruebas unitarias
4. Actualizar la documentación 