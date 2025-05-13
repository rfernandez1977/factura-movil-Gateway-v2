# Instrucciones para Ejecutar Pruebas en FMgo

## Problema Detectado

Se han identificado diversos problemas:

1. **Importaciones incorrectas**: El módulo está definido como `github.com/rodrigofernandezcalderon/FMgo` en el archivo go.mod, pero los archivos de código estaban importando paquetes como `fmgo/models`, lo que causa errores de compilación.

2. **Dependencia libxml2**: Se ha encontrado que el proyecto depende de la biblioteca `libxml2` a través del paquete Go `github.com/lestrrat-go/libxml2`. Esta dependencia requiere que `pkg-config` esté instalado en el sistema y configurado correctamente para encontrar libxml2.

3. **Puerto 8080 ocupado**: El servidor mock del SII estaba configurado para usar el puerto 8080, pero este puerto ya estaba siendo utilizado por otra aplicación.

4. **Funciones duplicadas**: Se detectaron funciones duplicadas en varios archivos de servicios, como `generateID()`.

5. **Estructuras XML incompletas**: Faltaban las estructuras de datos completas para generar sobres XML según las especificaciones del SII.

## Soluciones Implementadas

### 1. Configuración de Dependencias
- Se ha instalado `pkg-config` y `libxml2` mediante Homebrew:
  ```
  brew install pkg-config libxml2
  ```
- Se han establecido las variables de entorno necesarias:
  ```
  export PKG_CONFIG_PATH="/usr/local/opt/libxml2/lib/pkgconfig"
  export PATH="/usr/local/opt/libxml2/bin:$PATH"
  ```

### 2. Cambio de Puerto en el Servidor Mock
- Se ha modificado el puerto del servidor mock de 8080 a 9090 en el archivo `mock/sii_mock_server.go`.
- Se ha actualizado la referencia en la prueba `tests/emision_factura_test.go` para usar el puerto 9090.

### 3. Corrección de Importaciones
- Se han corregido las importaciones en todos los archivos para usar la ruta correcta `github.com/rodrigofernandezcalderon/FMgo`.

### 4. Desduplicación de Funciones
- Se ha centralizado la función `GenerateID()` en `services/utils_service.go`.
- Se han eliminado implementaciones duplicadas en `services/seguridad_service.go`, `services/ecommerce_service.go` y `services/reportes_service.go`.
- Se ha corregido la función duplicada `GetOrSet` en `services/cache_service.go`, manteniendo la versión más completa.

### 5. Implementación de Sobres XML
- Se han creado nuevos modelos basados en los esquemas XSD proporcionados por el SII:
  - `models/sobre_envio.go`: Implementa las estructuras para `EnvioDTE` y `EnvioBOLETA` según los esquemas `EnvioDTE_v10.xsd` y `EnvioBOLETA_v11.xsd`.
  - Soporte para los diferentes elementos XML requeridos, como `Caratula`, `SetDTE`, `Documento`, `TED`, etc.

- Se ha desarrollado un nuevo servicio `services/sobre_service.go` con funcionalidades para:
  - Generar sobres para facturas electrónicas
  - Generar sobres para boletas electrónicas
  - Crear sobres múltiples para el receptor
  - Crear sobres para el SII
  - Convertir los sobres a formato XML

- Se han implementado pruebas unitarias en `tests/sobre_service_test.go` para verificar la correcta generación de los sobres XML.

## Arquitectura de Sobres XML

En el sistema de facturación electrónica del SII, se requieren dos sobres:

1. **Sobre para el Receptor**: Contiene los documentos y va dirigido al cliente (RUT del receptor).
2. **Sobre para el SII**: Contiene los mismos documentos pero va dirigido al SII (RUT 60803000-K).

Cada sobre sigue esta estructura general:

```
<EnvioDTE>
  <SetDTE ID="...">
    <Caratula>
      <RutEmisor>...</RutEmisor>
      <RutEnvia>...</RutEnvia>
      <RutReceptor>...</RutReceptor>
      ...
    </Caratula>
    <DTE>
      <Documento>
        <Encabezado>...</Encabezado>
        <Detalle>...</Detalle>
        <TED>...</TED>
      </Documento>
      <Signature>...</Signature>
    </DTE>
    <!-- Más documentos si es necesario -->
  </SetDTE>
  <Signature>...</Signature>
</EnvioDTE>
```

## Pruebas Implementadas

1. **Simple Test**: `simple_test.go` - Verifica que el entorno Go funciona correctamente.
2. **Simple Factura Test**: `simple_factura_test.go` - Genera una factura XML básica sin dependencias complejas.
3. **Sobre Service Test**: `tests/sobre_service_test.go` - Prueba la generación de sobres XML para facturas, boletas y múltiples documentos.

## Instrucciones para Ejecutar Pruebas

Para ejecutar las pruebas simplificadas:

```
go test -v simple_test.go
go test -v simple_factura_test.go
```

Para probar la generación de sobres (requiere libxml2 y pkg-config):

```
cd tests
go test -v sobre_service_test.go
```

## Próximos Pasos

1. Implementar el proceso de firma de los sobres XML usando certificados digitales.
2. Finalizar la implementación de la integración con el servidor real del SII.
3. Agregar validación de los XML generados contra los esquemas XSD del SII.
4. Completar la información detallada de emisores y receptores (direcciones, comunas, etc.).
5. Implementar la integración completa del ciclo de emisión, envío y consulta de estado de documentos.