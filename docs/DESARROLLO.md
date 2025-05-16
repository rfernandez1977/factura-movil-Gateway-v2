# Guía de Desarrollo - FMgo

## Estándares de Código

### 1. Estructura de Código

#### 1.1 Organización de Directorios
```
.
├── core/              # Funcionalidad core del sistema
├── services/          # Servicios de la aplicación
├── models/           # Modelos de datos
├── middleware/       # Middleware HTTP
├── utils/           # Utilidades comunes
└── tests/           # Tests unitarios y de integración
```

#### 1.2 Nombrado
- Usar CamelCase para tipos exportados
- Usar snake_case para archivos
- Prefijos descriptivos para interfaces
- Sufijos descriptivos para implementaciones

```go
// Bien
type XMLProcessor interface {
    ProcessDocument(xml []byte) error
}

type DefaultXMLProcessor struct {
    // ...
}

// Mal
type Xml interface {
    process(data []byte) error
}
```

### 2. Documentación

#### 2.1 Comentarios
- Documentar todas las funciones exportadas
- Usar formato godoc
- Incluir ejemplos cuando sea relevante

```go
// FirmarXML firma un documento XML usando el certificado digital proporcionado.
// Retorna el documento firmado y un error si la operación falla.
//
// El documento debe estar en formato UTF-8 y ser válido según el schema XSD.
// 
// Ejemplo:
//
//     xmlData := []byte(`<Documento>...</Documento>`)
//     firmado, err := service.FirmarXML(xmlData)
func (s *FirmaService) FirmarXML(xmlData []byte) ([]byte, error) {
    // ...
}
```

#### 2.2 Documentación de Paquetes
- Incluir descripción general del paquete
- Documentar tipos principales
- Proporcionar ejemplos de uso

```go
// Package firma proporciona funcionalidades para la firma digital
// de documentos XML según los estándares del SII.
//
// Este paquete implementa XML-DSIG y maneja certificados digitales
// en formato P12/PEM.
package firma
```

### 3. Manejo de Errores

#### 3.1 Errores Personalizados
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validación fallida en %s: %s", e.Field, e.Message)
}
```

#### 3.2 Wrapping de Errores
```go
if err := validarXML(doc); err != nil {
    return fmt.Errorf("error validando documento: %w", err)
}
```

### 4. Testing

#### 4.1 Tests Unitarios
```go
func TestFirmarXML(t *testing.T) {
    tests := []struct {
        name    string
        input   []byte
        wantErr bool
    }{
        {
            name:    "documento válido",
            input:   []byte(`<Documento>...</Documento>`),
            wantErr: false,
        },
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewFirmaService()
            _, err := service.FirmarXML(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FirmarXML() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### 4.2 Tests de Integración
```go
func TestIntegracionSII(t *testing.T) {
    if testing.Short() {
        t.Skip("saltando test de integración")
    }
    // ...
}
```

### 5. Logging

#### 5.1 Niveles de Log
```go
// DEBUG para información de desarrollo
logger.Debug("Procesando documento", map[string]interface{}{
    "id": docID,
    "tipo": "factura",
})

// INFO para operaciones normales
logger.Info("Documento firmado exitosamente", map[string]interface{}{
    "id": docID,
})

// ERROR para errores que requieren atención
logger.Error("Error firmando documento", map[string]interface{}{
    "id": docID,
    "error": err.Error(),
})
```

#### 5.2 Contexto en Logs
```go
logger.Info("Operación completada", map[string]interface{}{
    "modulo": "firma",
    "operacion": "validacion",
    "duracion_ms": duration.Milliseconds(),
    "resultado": "exitoso",
})
```

### 6. Seguridad

#### 6.1 Manejo de Secretos
```go
// Bien
config := &Config{
    CertPassword: os.Getenv("CERT_PASSWORD"),
}

// Mal
config := &Config{
    CertPassword: "password123",
}
```

#### 6.2 Validación de Entrada
```go
func validarRUT(rut string) error {
    if !regexp.MatchString(`^\d{1,8}-[\dkK]$`, rut) {
        return fmt.Errorf("formato de RUT inválido")
    }
    // ...
}
```

### 7. Performance

#### 7.1 Uso de Buffer Pools
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func procesarXML(data []byte) error {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()
    // ...
}
```

#### 7.2 Optimización de Memoria
```go
// Bien
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    procesarLinea(scanner.Text())
}

// Mal
content, _ := ioutil.ReadAll(file)
lines := strings.Split(string(content), "\n")
for _, line := range lines {
    procesarLinea(line)
}
```

### 8. Concurrencia

#### 8.1 Goroutines Seguras
```go
func procesarLote(docs []Documento) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(docs))

    for _, doc := range docs {
        wg.Add(1)
        go func(d Documento) {
            defer wg.Done()
            if err := procesarDocumento(d); err != nil {
                errChan <- err
            }
        }(doc)
    }

    wg.Wait()
    close(errChan)

    return processErrors(errChan)
}
```

#### 8.2 Rate Limiting
```go
limiter := rate.NewLimiter(rate.Every(time.Second), 10)

func enviarAlSII(doc Documento) error {
    if err := limiter.Wait(context.Background()); err != nil {
        return err
    }
    // Enviar documento...
}
```

### 9. Configuración

#### 9.1 Variables de Entorno
```go
type Config struct {
    CertPath    string `envconfig:"CERT_PATH" required:"true"`
    KeyPath     string `envconfig:"KEY_PATH" required:"true"`
    SIIEndpoint string `envconfig:"SII_ENDPOINT" default:"https://palena.sii.cl"`
}
```

#### 9.2 Archivos de Configuración
```go
func LoadConfig(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("error leyendo config: %w", err)
    }

    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("error parseando config: %w", err)
    }

    return &config, nil
}
```

### 10. Versionado

#### 10.1 Semantic Versioning
- MAJOR.MINOR.PATCH
- MAJOR: cambios incompatibles
- MINOR: nuevas funcionalidades compatibles
- PATCH: correcciones de bugs

#### 10.2 Changelog
```markdown
## [1.2.1] - 2024-01-20
### Fixed
- Corrección en validación de certificados
- Mejora en manejo de memoria

## [1.2.0] - 2024-01-15
### Added
- Nuevo sistema de caché
- Soporte para certificados P12
```

## Flujo de Trabajo

### 1. Desarrollo
1. Crear rama feature/bugfix
2. Desarrollar con TDD
3. Documentar cambios
4. Crear tests
5. Actualizar documentación

### 2. Code Review
1. Self-review del código
2. Ejecutar linters
3. Verificar cobertura de tests
4. Solicitar revisión
5. Atender comentarios

### 3. Merge
1. Actualizar desde main
2. Resolver conflictos
3. Ejecutar tests
4. Merge a main
5. Tag version

## Herramientas Recomendadas

### 1. Desarrollo
- VS Code con extensiones Go
- Delve para debugging
- GoLand IDE

### 2. Linting
- golangci-lint
- staticcheck
- gosec

### 3. Testing
- go test
- testify
- gomock

### 4. Documentación
- godoc
- swag
- pkgsite

## Referencias

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Project Layout](https://github.com/golang-standards/project-layout) 