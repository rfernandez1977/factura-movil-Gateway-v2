# Guía de Certificación FMgo

## Ambiente de Certificación

### Configuración Inicial
1. **Certificados**
   - Ubicación: `./certificados/`
   - Archivos requeridos:
     - `cert_test.crt`: Certificado de prueba
     - `key_test.key`: Llave privada de prueba

2. **Esquemas XSD**
   - Ubicación: `./schemas/`
   - Archivos principales:
     - `SiiTypes_v10.xsd`: Tipos base
     - `DTE_v10.xsd`: Documentos tributarios
     - `EnvioDTE_v10.xsd`: Envío de documentos

3. **Configuración**
   - Archivo: `config/sii_config.json`
   - Parámetros principales:
     - Ambiente: certificación
     - URLs: maullin.sii.cl
     - Timeouts y reintentos
     - Monitoreo y métricas

### Monitoreo
1. **Logs**
   - Ubicación: `logs/certificacion/`
   - Niveles: DEBUG, INFO, WARN, ERROR
   - Rotación: diaria con compresión

2. **Métricas**
   - Ubicación: `metrics/certificacion/`
   - Métricas principales:
     - Tiempo de respuesta
     - Tasa de error
     - Uso de recursos

## Validaciones Implementadas

### 1. Validaciones de Documentos
- Formato XML según esquema XSD
- Montos y totales
- Referencias y relaciones
- Folios y secuencias

### 2. Validaciones de CAF
- Vigencia
- Rango de folios
- Firma del CAF
- Control de uso

### 3. Validaciones de Firma
- Certificado válido
- Firma XML-DSIG
- Integridad del documento
- Validez temporal

### 4. Reglas de Negocio
- Límites de montos
- Referencias requeridas
- Tipos de documentos válidos
- Secuencias y relaciones

## Procedimientos Operativos

### 1. Verificación de Ambiente
```bash
# Verificar configuración
cat config/sii_config.json

# Verificar certificados
ls -l certificados/

# Verificar esquemas
ls -l schemas/
```

### 2. Pruebas de Conectividad
```go
// Verificar conexión
client.VerificarComunicacion(ctx)

// Obtener semilla
semilla, err := client.ObtenerSemilla(ctx)

// Obtener token
token, err := client.ObtenerToken(ctx, semilla)
```

### 3. Monitoreo de Operación
```bash
# Verificar logs
tail -f logs/certificacion/fmgo.log

# Verificar métricas
cat metrics/certificacion/metrics.json
```

## Procedimientos de Contingencia

### 1. Errores de Conexión
1. Verificar conectividad con SII
2. Validar certificados
3. Revisar logs de error
4. Aplicar reintentos según configuración

### 2. Errores de Validación
1. Verificar esquemas XSD
2. Validar formato de documento
3. Revisar reglas de negocio
4. Verificar CAF y folios

### 3. Errores de Firma
1. Verificar certificado
2. Validar proceso de firma
3. Verificar timestamp
4. Revisar logs de firma

## Mantenimiento

### 1. Certificados
- Monitorear fechas de expiración
- Mantener respaldos seguros
- Actualizar cuando sea necesario

### 2. Esquemas
- Verificar actualizaciones del SII
- Mantener versiones anteriores
- Documentar cambios

### 3. Configuración
- Revisar periódicamente
- Ajustar según necesidad
- Mantener respaldos

## Referencias

1. [Documentación SII](https://www.sii.cl/factura_electronica/)
2. [Esquemas XSD](https://www.sii.cl/factura_electronica/factura_mercado/schema.html)
3. [Certificación DTE](https://www.sii.cl/factura_electronica/factura_mercado/proceso_certificacion.htm) 