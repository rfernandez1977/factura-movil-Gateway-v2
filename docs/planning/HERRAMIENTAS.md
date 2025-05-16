# Herramientas de Desarrollo - FMgo

## Herramientas Base

### Go Tools
```bash
# Versiones
go: 1.24.2
golangci-lint: v1.56.2
gotestsum: v1.11.0
mockgen: v1.6.0
swag: v1.16.3
```

### Herramientas de Testing
```bash
# Testing
go test
gotestsum
go-junit-report

# Cobertura
go tool cover
gocov

# Mocking
mockgen
gomock
```

### Herramientas de Análisis
```bash
# Linting
golangci-lint
staticcheck

# Seguridad
gosec
nancy

# Documentación
swag
godoc
```

### Herramientas de Base de Datos
```bash
# Migración
migrate
goose

# Testing
testcontainers-go
```

## Configuración de Desarrollo

### Editor/IDE
```bash
# VSCode Extensions
- Go
- Go Test Explorer
- Go Coverage
- Go Outline
- Go Doc
```

### Git Hooks
```bash
# Pre-commit
- gofmt
- golangci-lint
- go test
- gosec

# Pre-push
- go test ./...
- go test -race ./...
```

### Docker Compose
```yaml
services:
  postgres:
    image: postgres:16
    ports: ["5432:5432"]
  
  redis:
    image: redis:7
    ports: ["6379:6379"]
  
  mongodb:
    image: mongo:7
    ports: ["27017:27017"]
  
  rabbitmq:
    image: rabbitmq:3-management
    ports: ["5672:5672", "15672:15672"]
```

## Scripts de Desarrollo

### Instalación
```bash
scripts/
├── install-tools.sh
├── setup-dev.sh
├── setup-hooks.sh
└── setup-containers.sh
```

### Testing
```bash
scripts/
├── run-tests.sh
├── coverage.sh
└── integration-tests.sh
```

### Análisis
```bash
scripts/
├── lint.sh
├── security-check.sh
└── docs-generate.sh
```

## Configuración de CI/CD

### GitHub Actions
```yaml
workflows:
  - lint.yml
  - test.yml
  - security.yml
  - build.yml
  - deploy.yml
```

### Makefile Targets
```makefile
- setup
- test
- lint
- security
- build
- docs
```

## Métricas y Monitoreo

### Desarrollo
- Tiempo de compilación
- Cobertura de código
- Deuda técnica
- Vulnerabilidades

### Testing
- Tiempo de ejecución
- Fallos/Éxitos
- Cobertura
- Flaky tests

## Estándares de Código

### Formato
- gofmt
- goimports
- Configuración .editorconfig

### Linting
```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - golint
    - errcheck
    - staticcheck
    - gosec
```

### Testing
- Nombrado: Test{Func}__{Scenario}
- Tabla de casos
- Setup/Teardown
- Mocks consistentes

## Documentación

### Código
- Comentarios godoc
- Ejemplos de uso
- Tests como documentación

### API
- Swagger/OpenAPI
- Postman collections
- Ejemplos de requests

## Mantenimiento

### Diario
- Actualizar dependencias
- Ejecutar tests
- Verificar linting

### Semanal
- Análisis de seguridad
- Revisión de cobertura
- Actualización de docs

### Mensual
- Auditoría de dependencias
- Revisión de métricas
- Actualización de herramientas 