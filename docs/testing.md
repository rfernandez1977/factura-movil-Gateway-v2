# Documentación de Pruebas

## Estructura de Pruebas

```
tests/
├── unit/           # Pruebas unitarias
├── integration/    # Pruebas de integración
├── e2e/           # Pruebas end-to-end
└── fixtures/      # Datos de prueba
```

## Pruebas Unitarias

### Servicios

#### SIIClient
```go
func TestSIIClient_ObtenerSemilla(t *testing.T) {
    // Arrange
    client := NewMockSIIClient()
    
    // Act
    semilla, err := client.ObtenerSemilla()
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, semilla)
}
```

#### CAFManager
```go
func TestCAFManager_ValidarCAF(t *testing.T) {
    // Arrange
    manager := NewCAFManager(config)
    caf := &models.CAFXML{
        DA: struct {
            RNG struct {
                D int `xml:"D"`
                H int `xml:"H"`
            } `xml:"RNG"`
        }{
            RNG: struct {
                D int `xml:"D"`
                H int `xml:"H"`
            }{
                D: 1,
                H: 100,
            },
        },
    }
    
    // Act
    err := manager.ValidarCAF(caf, 50)
    
    // Assert
    assert.NoError(t, err)
}
```

#### DTEGenerator
```go
func TestDTEGenerator_GenerarDTE(t *testing.T) {
    // Arrange
    generator := NewDTEGenerator(caf)
    emisor := models.Emisor{
        RUT: "76.123.456-7",
        RazonSocial: "Empresa de Prueba",
    }
    receptor := models.Receptor{
        RUT: "76.765.432-1",
        RazonSocial: "Cliente de Prueba",
    }
    items := []models.Item{
        {
            Descripcion: "Producto 1",
            Cantidad: 1,
            PrecioUnitario: 1000,
        },
    }
    
    // Act
    dte, err := generator.GenerarDTE(emisor, receptor, items)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, dte)
    assert.Equal(t, 1000.0, dte.MontoNeto)
    assert.Equal(t, 190.0, dte.MontoIVA)
    assert.Equal(t, 1190.0, dte.MontoTotal)
}
```

## Pruebas de Integración

### Flujo Completo
```go
func TestFlujoCompleto(t *testing.T) {
    // Arrange
    client := NewSIIClient(config)
    manager := NewCAFManager(config)
    generator := NewDTEGenerator(caf)
    
    // Act
    // 1. Obtener semilla
    semilla, err := client.ObtenerSemilla()
    assert.NoError(t, err)
    
    // 2. Obtener token
    token, err := client.ObtenerToken(semilla)
    assert.NoError(t, err)
    
    // 3. Generar DTE
    dte, err := generator.GenerarDTE(emisor, receptor, items)
    assert.NoError(t, err)
    
    // 4. Enviar al SII
    err = client.EnviarDTE(dte, token)
    assert.NoError(t, err)
    
    // 5. Consultar estado
    estado, err := client.ConsultarEstado(dte.TrackID)
    assert.NoError(t, err)
    assert.Equal(t, "ACEPTADO", estado.Estado)
}
```

## Pruebas End-to-End

### Emisión de Factura
```go
func TestEmisionFactura(t *testing.T) {
    // Arrange
    server := NewTestServer()
    defer server.Close()
    
    // Act
    response, err := server.EmitirFactura(factura)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, response.StatusCode)
    
    var result models.DocumentoTributario
    err = json.NewDecoder(response.Body).Decode(&result)
    assert.NoError(t, err)
    assert.Equal(t, "ACEPTADO", result.Estado)
}
```

## Datos de Prueba

### Fixtures
```json
{
  "empresa": {
    "rut": "76.123.456-7",
    "razon_social": "Empresa de Prueba",
    "giro": "Servicios de Prueba",
    "direccion": "Calle 123",
    "comuna": "Santiago",
    "ciudad": "Santiago"
  },
  "documentos": [
    {
      "tipo": "FACTURA",
      "folio": 1,
      "fecha_emision": "2024-03-20",
      "monto_neto": 1000,
      "monto_iva": 190,
      "monto_total": 1190,
      "items": [
        {
          "descripcion": "Producto 1",
          "cantidad": 1,
          "precio_unitario": 1000
        }
      ]
    }
  ]
}
```

## Mock Server

### Configuración
```yaml
port: 8081
endpoints:
  - path: /DTEWS/CrSeed.jws
    method: POST
    response:
      status: 200
      body: ./fixtures/semilla_response.xml
  - path: /DTEWS/GetTokenFromSeed.jws
    method: POST
    response:
      status: 200
      body: ./fixtures/token_response.xml
  - path: /cgi_dte/UPL/DTEUpload
    method: POST
    response:
      status: 200
      body: ./fixtures/envio_response.xml
```

### Respuestas Mock
```xml
<!-- fixtures/semilla_response.xml -->
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <getSeedResponse>
      <seed>SEMILLA_TEST</seed>
    </getSeedResponse>
  </soap:Body>
</soap:Envelope>
```

## Cobertura de Pruebas

### Métricas
- Cobertura de código: >80%
- Pruebas unitarias: >70%
- Pruebas de integración: >20%
- Pruebas e2e: >10%

### Reportes
```bash
# Generar reporte de cobertura
go test -coverprofile=coverage.out ./...

# Ver reporte en HTML
go tool cover -html=coverage.out
```

## Ejecución de Pruebas

### Comandos
```bash
# Ejecutar todas las pruebas
make test

# Ejecutar pruebas unitarias
make test-unit

# Ejecutar pruebas de integración
make test-integration

# Ejecutar pruebas e2e
make test-e2e

# Ejecutar pruebas con cobertura
make test-coverage
```

### CI/CD
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run tests
        run: make test
      - name: Upload coverage
        uses: codecov/codecov-action@v1
```

## Mejores Prácticas

### Naming
- Usar nombres descriptivos
- Seguir el patrón `Test<Componente>_<Acción>`
- Incluir el caso de prueba en el nombre

### Estructura
- Seguir el patrón AAA (Arrange, Act, Assert)
- Mantener las pruebas independientes
- Usar fixtures para datos de prueba

### Aserciones
- Usar aserciones específicas
- Verificar el estado final
- Validar mensajes de error

### Limpieza
- Limpiar recursos después de las pruebas
- Usar `t.Cleanup()`
- Implementar `TestMain(m *testing.M)` 