# Documentación de Pruebas - Cliente SII

## Estructura de Pruebas

El cliente SII está compuesto por varios componentes que trabajan juntos para proporcionar la funcionalidad completa de comunicación con el SII. Las pruebas están organizadas según estos componentes:

### 1. Cliente SOAP (`soap_client_test.go`)
- `TestNewSOAPClient`: Valida la creación correcta del cliente SOAP
- `TestSOAPClientCall`: Prueba las llamadas SOAP básicas
- `TestSOAPClientErrors`: Verifica el manejo de errores en llamadas SOAP

### 2. Cliente de Autenticación (`auth_client_test.go`)
- `TestNewAuthClient`: Valida la creación del cliente de autenticación
- `TestGetToken`: Prueba la obtención de tokens
- `TestGetTokenErrors`: Verifica el manejo de errores en la autenticación
- `TestConcurrentTokenAccess`: Valida el acceso concurrente a tokens

### 3. Cliente DTE (`dte_client_test.go`)
- `TestNewDTEClient`: Valida la creación del cliente DTE
- `TestEnviarDTE`: Prueba el envío de documentos tributarios
- `TestConsultarEstadoDTE`: Verifica la consulta de estado de documentos
- `TestConsultarEstadoEnvio`: Prueba la consulta de estado de envíos
- `TestDTEClientErrors`: Valida el manejo de errores
- `TestDTEClientValidations`: Verifica las validaciones de datos
- `TestDTEClientTimeout`: Prueba el manejo de timeouts
- `TestDTEClientConcurrency`: Valida operaciones concurrentes
- `TestDTEClientRetry`: Prueba la política de reintentos
- `TestDTEClientContextCancellation`: Verifica la cancelación de contexto

### 4. Cliente HTTP (`http_client_test.go`)
- `TestHTTPClient_EnviarDTE`: Prueba el envío HTTP de documentos
- `TestHTTPClient_ConsultarEstado`: Verifica consultas HTTP de estado
- `TestHTTPClient_ConsultarDTE`: Prueba consultas HTTP de documentos

### 5. Manejo de Certificados (`certificates_test.go`)
- Pruebas de carga y validación de certificados digitales
- Verificación de firmas y encriptación

### 6. Sistema de Reintentos (`retry_test.go`)
- `TestWithRetry_Success`: Prueba reintentos exitosos
- `TestWithRetry_EventualSuccess`: Verifica éxito después de reintentos
- `TestWithRetry_Failure`: Valida fallas después de reintentos máximos
- `TestWithRetry_ContextCancellation`: Prueba cancelación durante reintentos

## Configuración de Pruebas

### Archivos de Prueba
Los archivos necesarios para las pruebas se encuentran en:
```
test/
├── fixtures/          # Datos de prueba
│   ├── certs/        # Certificados de prueba
│   └── responses/    # Respuestas XML de ejemplo
└── mocks/            # Implementaciones mock
```

### Certificados de Prueba
Se utilizan certificados de prueba generados específicamente para testing:
- Certificado de firma digital (.p12)
- Certificado SSL para conexiones HTTPS
- Llaves privadas y públicas para pruebas

## Casos de Prueba Principales

### 1. Envío de DTE
```go
// TestEnviarDTE verifica:
// - Autenticación correcta
// - Formato XML válido
// - Firma digital correcta
// - Respuesta exitosa del SII
// - Manejo de trackID
```

### 2. Consulta de Estado
```go
// TestConsultarEstadoDTE verifica:
// - Token válido
// - Parámetros correctos
// - Parsing de respuesta
// - Estados diferentes
```

### 3. Manejo de Errores
```go
// TestDTEClientErrors verifica:
// - Errores de validación
// - Errores de conexión
// - Errores de timeout
// - Errores de formato
```

### 4. Concurrencia
```go
// TestDTEClientConcurrency verifica:
// - Envíos simultáneos
// - Manejo de recursos
// - Race conditions
```

## Mejores Prácticas

### 1. Preparación de Pruebas
- Usar `setupTestFiles()` para crear archivos temporales
- Limpiar recursos con `defer cleanup()`
- Configurar timeouts apropiados

### 2. Assertions
- Usar `require` para condiciones críticas
- Usar `assert` para validaciones no críticas
- Verificar errores específicos

### 3. Mocking
- Implementar interfaces para facilitar mocking
- Usar servidores de prueba para simular SII
- Simular latencias y errores

### 4. Documentación
- Documentar propósito de cada prueba
- Incluir ejemplos de uso
- Explicar casos edge

## Mantenimiento

### 1. Cobertura
- Mantener cobertura >80%
- Verificar casos edge
- Actualizar pruebas al modificar código

### 2. Performance
- Monitorear tiempo de ejecución
- Optimizar pruebas lentas
- Paralelizar cuando sea posible

### 3. Actualizaciones
- Revisar pruebas periódicamente
- Actualizar según cambios del SII
- Mantener datos de prueba actualizados 