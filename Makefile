.PHONY: build run test clean

# Variables
BINARY_NAME=api
BUILD_DIR=build
CMD_DIR=cmd/api

# Comandos
build:
	@echo "Compilando..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

run: build
	@echo "Ejecutando..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

test:
	@echo "Ejecutando tests..."
	@go test -v ./...

clean:
	@echo "Limpiando..."
	@rm -rf $(BUILD_DIR)

# Desarrollo
dev:
	@echo "Iniciando en modo desarrollo..."
	@go run ./$(CMD_DIR)

# Dependencias
deps:
	@echo "Instalando dependencias..."
	@go mod download
	@go mod tidy

# Linting
lint:
	@echo "Ejecutando linter..."
	@golangci-lint run

# Formateo
fmt:
	@echo "Formateando código..."
	@go fmt ./...

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  make build    - Compila el proyecto"
	@echo "  make run      - Ejecuta el proyecto"
	@echo "  make test     - Ejecuta los tests"
	@echo "  make clean    - Limpia los archivos compilados"
	@echo "  make dev      - Ejecuta en modo desarrollo"
	@echo "  make deps     - Instala dependencias"
	@echo "  make lint     - Ejecuta el linter"
	@echo "  make fmt      - Formatea el código" 