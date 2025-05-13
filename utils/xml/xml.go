package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// MarshalToFile convierte una estructura a XML y la guarda en un archivo
func MarshalToFile(v interface{}, filename string) error {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("error al convertir a XML: %w", err)
	}

	// Agregar la declaración XML
	xmlData := []byte(xml.Header + string(data))

	err = os.WriteFile(filename, xmlData, 0644)
	if err != nil {
		return fmt.Errorf("error al escribir archivo: %w", err)
	}

	return nil
}

// UnmarshalFromFile lee un archivo XML y lo convierte a una estructura
func UnmarshalFromFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error al leer archivo: %w", err)
	}

	err = xml.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("error al convertir XML: %w", err)
	}

	return nil
}

// MarshalToString convierte una estructura a XML y retorna el string
func MarshalToString(v interface{}) (string, error) {
	data, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al convertir a XML: %w", err)
	}

	return xml.Header + string(data), nil
}

// UnmarshalFromString convierte un string XML a una estructura
func UnmarshalFromString(data string, v interface{}) error {
	err := xml.Unmarshal([]byte(data), v)
	if err != nil {
		return fmt.Errorf("error al convertir XML: %w", err)
	}

	return nil
}

// MarshalToWriter convierte una estructura a XML y la escribe en un io.Writer
func MarshalToWriter(v interface{}, w io.Writer) error {
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	// Escribir la declaración XML
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return fmt.Errorf("error al escribir declaración XML: %w", err)
	}

	err = encoder.Encode(v)
	if err != nil {
		return fmt.Errorf("error al codificar XML: %w", err)
	}

	return nil
}

// UnmarshalFromReader lee XML de un io.Reader y lo convierte a una estructura
func UnmarshalFromReader(r io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("error al decodificar XML: %w", err)
	}

	return nil
}

// ValidateXML valida que un string sea XML válido
func ValidateXML(data string) error {
	decoder := xml.NewDecoder(strings.NewReader(data))
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML inválido: %w", err)
		}
	}
	return nil
}

// PrettyPrintXML formatea un XML para que sea más legible
func PrettyPrintXML(data string) (string, error) {
	var v interface{}
	err := UnmarshalFromString(data, &v)
	if err != nil {
		return "", err
	}

	pretty, err := MarshalToString(v)
	if err != nil {
		return "", err
	}

	return pretty, nil
}
