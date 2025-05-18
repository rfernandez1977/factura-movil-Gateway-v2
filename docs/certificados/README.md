# Configuración de Certificados Digitales

## Estructura de Directorios

```
dev/
└── config/
    └── certs/
        └── firma/
            ├── config.json
            ├── firma.key
            └── firmaFM.pfx
```

## Archivos de Certificado

### 1. Certificado Digital (firmaFM.pfx)
- **Ubicación**: `dev/config/certs/firma/firmaFM.pfx`
- **Tipo**: Certificado PKCS#12
- **Uso**: Contiene el certificado digital y la clave privada
- **Seguridad**: Protegido con contraseña

### 2. Llave Privada (firma.key)
- **Ubicación**: `dev/config/certs/firma/firma.key`
- **Tipo**: Archivo de llave privada
- **Uso**: Utilizado para firmar documentos digitalmente

### 3. Configuración (config.json)
- **Ubicación**: `dev/config/certs/firma/config.json`
- **Contenido**:
  ```json
  {
      "certificado": {
          "rut_firmante": "13195458-1",
          "ruta_key": "dev/config/certs/firma/firma.key",
          "ruta_pfx": "dev/config/certs/firma/firmaFM.pfx",
          "password": "83559705FM"
      }
  }
  ```

## Validaciones Implementadas

Se han implementado las siguientes validaciones en `core/firma/test/firma_test.go`:

1. Verificación de existencia de archivos
   - Certificado digital (.pfx)
   - Llave privada (.key)
   - Archivo de configuración

2. Validación del formato RUT
   - Formato: `^\d{1,8}-[\dkK]$`
   - Ejemplo válido: "13195458-1"

3. Validación de estructura de configuración
   - Parseo correcto del JSON
   - Campos requeridos presentes

## Seguridad

### Permisos de Archivos
- **config.json**: 644 (rw-r--r--)
- **firma.key**: 644 (rw-r--r--)
- **firmaFM.pfx**: 644 (rw-r--r--)

### Recomendaciones de Seguridad
1. No versionar los certificados en Git
2. Mantener respaldos seguros
3. Rotar contraseñas periódicamente
4. Monitorear accesos a los archivos

## Uso en Ambiente de Desarrollo

Para utilizar los certificados en desarrollo:

1. Copiar los archivos al directorio correspondiente:
   ```bash
   mkdir -p dev/config/certs/firma
   cp firma.key dev/config/certs/firma/
   cp firmaFM.pfx dev/config/certs/firma/
   ```

2. Verificar la configuración:
   ```bash
   cd core/firma/test
   go test -v ./... -run TestCertificadoDigital
   ```

## Próximos Pasos

1. Implementación del servicio de firma
2. Pruebas de integración con SII
3. Configuración de ambiente de certificación
4. Implementación de logs de auditoría 