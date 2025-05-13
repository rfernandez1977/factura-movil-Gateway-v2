# FMgo - Sistema de Facturación Electrónica

Este proyecto implementa un sistema de facturación electrónica compatible con el SII (Servicio de Impuestos Internos) de Chile.

## Cambios recientes en la estructura de modelos

- **Modelos de negocio centralizados:** Todos los modelos principales (Factura, Boleta, NotaCredito, NotaDebito, GuiaDespacho, DocumentoTributario, Item, etc.) están centralizados en el directorio `models/`.
- **Modelos auxiliares para interoperabilidad:** Si se requiere interoperar con servicios externos (por ejemplo, generación/parsing de XML para el SII), existen modelos auxiliares como `DTEDocument` en `models/` o `ItemXML` en `utils/`.
- **No duplicar modelos:** No se deben duplicar modelos en `utils/` o en archivos de test. Los tests deben importar y usar los modelos reales desde `models/`.
- **Stubs y mocks:** Los archivos mock y stub solo deben usarse para pruebas y nunca en producción.

## Patrón de uso de modelos

- **Para lógica de negocio interna:** Usa siempre los modelos de `models/`.
- **Para interoperabilidad (XML, SII, etc.):** Usa los modelos auxiliares (`DTEDocument`, `ItemXML`, etc.) y documenta claramente su propósito.
- **Para pruebas:** Usa los modelos reales de `models/` y mocks solo para simular dependencias externas.

## Estructura del proyecto

- `models/`: Definición de estructuras de datos principales y auxiliares para interoperabilidad.
- `services/`: Implementación de servicios (SII, XML, etc.).
- `utils/`: Utilidades y herramientas auxiliares. Aquí solo deben existir modelos auxiliares para interoperabilidad, nunca duplicados de negocio.
- `tests/`: Pruebas del sistema. Deben importar modelos desde `models/`.
- `mock/`: Implementación del servidor mock para pruebas.
- `config/`: Archivos de configuración.
- `test_cases/`: Casos de prueba.

## Recomendaciones para desarrolladores

- Antes de crear un nuevo modelo, revisa si ya existe en `models/`.
- Si necesitas un modelo para interoperabilidad, documenta claramente su propósito y nómbralo diferente (ej: `DTEDocument`, `ItemXML`).
- No dupliques modelos en utilidades o tests.
- Si encuentras un stub o placeholder, reemplázalo por una implementación real o documenta claramente que es solo para pruebas.

## Pruebas con Mock del SII

Para facilitar las pruebas sin necesidad de conectarse a los servicios reales del SII, se ha implementado un servidor mock que simula las respuestas del SII.

### Configuración para pruebas

1. **Configuración de empresa demo**: 
   - La configuración de la empresa de prueba (Factura Movil SPA) está en `config/config_demo_facturamovil.json`

2. **Casos de prueba**:
   - Casos de prueba disponibles en `test_cases/`
   - El caso de emisión de factura está en `test_cases/emision_factura_test.json`

3. **CAF y firma digital**:
   - CAF (Código de Autorización de Folios) para pruebas en `caf_test/33.xml`
   - Certificado de firma digital simulado en `firma_test/firma.p12`

### Ejecución de pruebas

Para ejecutar las pruebas, se necesitan dos terminales:

#### Terminal 1: Iniciar el servidor mock del SII

```bash
make mock-server
```

Este comando inicia un servidor HTTP en el puerto 8080 que simula las respuestas del SII.

#### Terminal 2: Ejecutar las pruebas

```bash
make test-emision
```

Este comando ejecuta la prueba de emisión de factura, que:
1. Genera una factura basada en el caso de prueba
2. Genera el XML correspondiente
3. Simula el envío al SII (servidor mock)
4. Consulta el estado del documento usando el TrackID recibido

### Otros comandos disponibles

```bash
# Ejecutar todos los tests
make test

# Limpiar archivos temporales
make clean

# Mostrar ayuda
make help
```

## Requisitos

- Go 1.18 o superior
- Dependencias gestionadas con Go Modules

## Descripción
FMgo es un sistema de facturación electrónica que integra múltiples plataformas de e-commerce y se conecta con el Servicio de Impuestos Internos (SII) de Chile para la emisión de documentos tributarios electrónicos.

## Características Principales
- Integración con SII para emisión de documentos electrónicos
- Sincronización con múltiples plataformas de e-commerce:
  - Shopify
  - PrestaShop
  - WooCommerce
  - Jumpseller
- Dashboard de administración con Grafana
- Sistema de monitoreo y métricas
- API RESTful para integración
- Generación de PDFs
- Sistema de alertas y notificaciones

## Instalación

### 1. Clonar el repositorio
```bash
git clone https://github.com/tu-usuario/fmgo.git
cd fmgo
```

### 2. Instalar dependencias
```bash
go mod download
```

### 3. Configurar variables de entorno
Crear un archivo `.env` con las siguientes variables:
```env
SII_BASE_URL=https://palena.sii.cl
CERT_PATH=./certs/cert.pem
KEY_PATH=./certs/key.pem
SII_AMBIENTE=CERTIFICACION
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=fmgo
```

### 4. Iniciar el servidor
```bash
go run main.go
```

## Estructura del Proyecto
```
fmgo/
├── api/              # Definiciones de API
├── config/           # Archivos de configuración
├── controllers/      # Controladores
├── db/              # Configuración de base de datos
├── docs/            # Documentación
├── handlers/        # Manejadores de peticiones
├── middleware/      # Middleware
├── migrations/      # Migraciones de base de datos
├── models/          # Modelos de datos
├── repository/      # Repositorios
├── routes/          # Rutas
├── services/        # Servicios de negocio
├── utils/           # Utilidades
└── metrics/         # Métricas y monitoreo
```

## API Documentation
La documentación detallada de la API se encuentra en [docs/gateway-api-endpoints.md](docs/gateway-api-endpoints.md).

## Dashboard de Administración
El dashboard de administración está configurado en Grafana y se encuentra en `config/company-admin-dashboard.json`. Incluye:
- Resumen de compañías
- Configuraciones por plataforma
- Estado de sincronización
- Estadísticas de uso de API
- Configuración de webhooks
- Sistema de alertas
- Logs de auditoría

## Monitoreo
El sistema utiliza Prometheus para el monitoreo de métricas. Las métricas están disponibles en el endpoint `/metrics`.

## Contribución
1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## Licencia
Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## Contacto
Para soporte o consultas, contactar a [tu-email@ejemplo.com](mailto:tu-email@ejemplo.com)

## Arquitectura

El sistema sigue una arquitectura de capas:

1. **Modelos**: Representan las estructuras de datos utilizadas en la aplicación.
2. **Servicios**: Contienen la lógica de negocio para procesar documentos tributarios.
3. **Controladores**: Manejan las solicitudes HTTP y coordinan los servicios.
4. **Repositorios**: Gestionan la persistencia de datos.

## Componentes Principales

### Sobres XML

El sistema implementa dos tipos de sobres XML según los requerimientos del SII:

1. **Sobre para el Receptor (Cliente)**: Contiene los documentos tributarios enviados al cliente.
2. **Sobre para el SII**: Contiene los documentos tributarios que deben ser informados al SII.

La estructura de estos sobres está definida en los esquemas XSD proporcionados por el SII:

- `EnvioDTE_v10.xsd`: Define la estructura del sobre para facturas.
- `EnvioBOLETA_v11.xsd`: Define la estructura del sobre para boletas.

### Firma Digital

Para cumplir con los requisitos del SII, el sistema implementa la firma digital de los sobres XML utilizando el estándar XML-DSIG (XML Signature). La firma se aplica a diferentes elementos del sobre:

1. **Firma del Sobre**: Se firma el sobre completo para garantizar su integridad.
2. **Firma de cada Documento**: Cada documento dentro del sobre (factura o boleta) también se firma individualmente.

La implementación de firma digital incluye:

- **FirmaDigitalService**: Servicio responsable de firmar documentos XML utilizando certificados digitales.
- **Algoritmos**: Se utiliza RSA con SHA-1 según los requerimientos del SII.
- **Certificados**: Se utilizan certificados digitales emitidos por entidades autorizadas.

Para utilizar la firma digital, es necesario:

1. Obtener un certificado digital válido desde una entidad autorizada por el SII.
2. Configurar el servicio de firma con la ruta al certificado y la llave privada.
3. Integrar el servicio de firma con el servicio de sobres.

Ejemplo de uso:

```go
// Crear servicio de firma digital
firmaService, err := services.NewFirmaDigitalService("ruta/al/certificado.crt", "ruta/a/llave.key", "contraseña", "11111111-1")
if err != nil {
    log.Fatalf("Error al crear servicio de firma: %v", err)
}

// Crear servicio de sobres con firma digital
sobreService := services.NewSobreService(firmaService)

// Generar un sobre firmado
sobre, err := sobreService.GenerarSobreFactura(factura, fechaResolucion, numeroResolucion)
if err != nil {
    log.Fatalf("Error al generar sobre: %v", err)
}

// Convertir a XML (incluye la firma digital)
xmlString, err := sobreService.ConvertirAXML(sobre)
if err != nil {
    log.Fatalf("Error al convertir a XML: %v", err)
}
```

## Próximos Pasos

- Implementar envío de documentos al SII a través de sus APIs
- Mejorar manejo de errores y validaciones
- Agregar más pruebas de integración 