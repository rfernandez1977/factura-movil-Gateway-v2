package services

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidator valida esquemas JSON
type SchemaValidator struct {
	schemas map[string]interface{}
}

// NewSchemaValidator crea un nuevo validador de esquemas
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemas: make(map[string]interface{}),
	}
}

// RegisterSchema registra un esquema para validación
func (v *SchemaValidator) RegisterSchema(name string, schema interface{}) error {
	if _, exists := v.schemas[name]; exists {
		return fmt.Errorf("el esquema %s ya está registrado", name)
	}
	v.schemas[name] = schema
	return nil
}

// ValidateJSON valida un JSON contra un esquema registrado
func (v *SchemaValidator) ValidateJSON(name string, data []byte) error {
	schema, exists := v.schemas[name]
	if !exists {
		return fmt.Errorf("esquema %s no encontrado", name)
	}

	// Convertir el esquema a JSON
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("error al convertir esquema a JSON: %v", err)
	}

	// Crear el validador
	schemaLoader := gojsonschema.NewBytesLoader(schemaJSON)
	documentLoader := gojsonschema.NewBytesLoader(data)

	// Validar
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error al validar JSON: %v", err)
	}

	if !result.Valid() {
		var errors string
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("JSON no válido:\n%s", errors)
	}

	return nil
}
