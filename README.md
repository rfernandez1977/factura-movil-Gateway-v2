# FMgo - Sistema de Integración SII

Sistema de integración con el Servicio de Impuestos Internos (SII) de Chile para la gestión de documentos tributarios electrónicos.

## Características Principales

- Validación XSD de documentos tributarios
- Sistema de firma digital avanzado con soporte XML-DSIG
- Envío automático al SII con reintentos y manejo de errores
- Gestión completa de DTE (Documentos Tributarios Electrónicos)
- Sistema de monitoreo y logging multinivel
- Caché de certificados digitales con rotación automática
- Validación de firmas y certificados
- Soporte para múltiples formatos de certificados (P12/PEM)

## Requisitos

- Go 1.23 o superior
- libxml2
- PostgreSQL 14+
- Redis 7+
- OpenSSL
- Certificado digital válido del SII

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/FMgo.git

# Instalar dependencias
go mod download

# Configurar entorno
cp .env.example .env
cp config.example.json config.json

# Compilar
make build
```

## Configuración

1. Configurar variables de entorno en `.env`:
   ```env
   CERT_PATH=/ruta/al/certificado.p12
   KEY_PATH=/ruta/a/llave.key
   CERT_PASSWORD=tu_password
   RUT_EMPRESA=76.555.555-5
   ```

2. Configurar `config.json` con los parámetros de conexión
3. Asegurar que los esquemas XSD estén en `schema_dte/`
4. Configurar los niveles de log en `logging.json`

## Documentación

### Guías Técnicas
- [Arquitectura del Sistema](docs/ARQUITECTURA.md)
- [Sistema de Firma Digital](docs/FIRMA_DIGITAL.md)
- [Sistema de Logging](docs/LOGGING.md)
- [API Reference](docs/API.md)
- [Guía de Desarrollo](docs/DESARROLLO.md)
- [Troubleshooting](docs/TROUBLESHOOTING.md)

### Ejemplos de Uso

#### Sistema de Firma Digital
```go
// Crear servicio de firma
firmaService, err := services.NewXMLSignatureService(
    os.Getenv("CERT_PATH"),
    os.Getenv("KEY_PATH"),
    os.Getenv("CERT_PASSWORD"),
    os.Getenv("RUT_EMPRESA")
)

// Firmar documento
signedXML, err := firmaService.FirmarXML(xmlData)

// Validar firma
isValid, err := firmaService.ValidarFirma(signedXML)
```

#### Sistema de Logging
```go
// Inicializar logger con rotación
logger, err := logger.NewLogger("logs/app.log", logger.DEBUG)
defer logger.Close()

// Logging estructurado
logger.Info("Proceso iniciado", map[string]interface{}{
    "modulo": "firma",
    "operacion": "validacion",
    "documento_id": "DOC001"
})

// Logging de operaciones XML
logger.LogXMLOperation("FirmarXML", xmlData, err)
```

## Estructura del Proyecto

```
.
├── docs/               # Documentación
├── core/              # Núcleo del sistema
│   ├── sii/          # Integración con SII
│   ├── firma/        # Servicios de firma
│   └── logger/       # Sistema de logging
├── models/           # Modelos de datos
├── services/         # Servicios principales
├── middleware/       # Middleware
├── utils/           # Utilidades comunes
├── tests/           # Tests
│   ├── unit/        # Tests unitarios
│   └── integration/ # Tests de integración
├── schema_dte/      # Esquemas XSD
└── scripts/         # Scripts de utilidad
```

## Métricas y Monitoreo

- Cobertura de tests: > 80%
- Tiempo de respuesta API: < 200ms
- Uptime: > 99.9%
- Monitoreo en tiempo real vía Prometheus/Grafana
- Alertas configurables por nivel de severidad

## Contribución

1. Fork el proyecto
2. Crear rama feature (`git checkout -b feature/NuevaFuncionalidad`)
3. Commit cambios (`git commit -m 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/NuevaFuncionalidad`)
5. Crear Pull Request

### Guías de Contribución

- Seguir estándares de código Go
- Documentar nuevas funcionalidades
- Mantener cobertura de tests > 80%
- Usar logging estructurado
- Seguir principios SOLID

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Soporte

Para soporte técnico:
- Email: soporte@fmgo.cl
- Documentación: [https://docs.fmgo.cl](https://docs.fmgo.cl)
- Issues: [https://github.com/tu-usuario/FMgo/issues](https://github.com/tu-usuario/FMgo/issues) 