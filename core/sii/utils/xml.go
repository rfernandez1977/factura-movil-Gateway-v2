package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// GenerateXML genera un XML a partir de una estructura
func GenerateXML(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(xml.Header)

	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("error al generar XML: %v", err)
	}

	return buf.Bytes(), nil
}

// ParseXML parsea un XML en una estructura
func ParseXML(data []byte, v interface{}) error {
	if err := xml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("error al parsear XML: %v", err)
	}
	return nil
}

// ValidateXML valida que un XML esté bien formado
func ValidateXML(data []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return fmt.Errorf("error al validar XML: %v", err)
		}
	}
}
