package xmlutils

import (
	"os"
	"testing"
)

type TestStruct struct {
	XMLName struct{} `xml:"TestStruct"`
	Name    string   `xml:"name"`
	Value   int      `xml:"value"`
}

type TestSOAPResponse struct {
	XMLName struct{} `xml:"Envelope"`
	Body    struct {
		Response struct {
			Result string `xml:"result"`
		} `xml:"TestResponse"`
	} `xml:"Body"`
}

func TestXMLParser_ParseXML(t *testing.T) {
	parser := NewXMLParser(false)

	xmlData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<TestStruct>
    <name>Test</name>
    <value>123</value>
</TestStruct>`)

	var result TestStruct
	err := parser.ParseXML(xmlData, &result)
	if err != nil {
		t.Errorf("ParseXML failed: %v", err)
	}

	if result.Name != "Test" {
		t.Errorf("Expected name 'Test', got '%s'", result.Name)
	}
	if result.Value != 123 {
		t.Errorf("Expected value 123, got %d", result.Value)
	}
}

func TestXMLParser_ParseSOAP(t *testing.T) {
	parser := NewXMLParser(false)

	soapData := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <ns1:TestResponse xmlns:ns1="http://test.com">
            <result>Success</result>
        </ns1:TestResponse>
    </soap:Body>
</soap:Envelope>`)

	var result TestSOAPResponse
	err := parser.ParseSOAP(soapData, &result)
	if err != nil {
		t.Errorf("ParseSOAP failed: %v", err)
	}

	if result.Body.Response.Result != "Success" {
		t.Errorf("Expected result 'Success', got '%s'", result.Body.Response.Result)
	}
}

func TestXMLParser_GenerateXML(t *testing.T) {
	parser := NewXMLParser(false)

	testData := TestStruct{
		Name:  "Test",
		Value: 123,
	}

	xmlData, err := parser.GenerateXML(testData)
	if err != nil {
		t.Errorf("GenerateXML failed: %v", err)
	}

	// Verificar que podemos parsear el XML generado
	var result TestStruct
	err = parser.ParseXML(xmlData, &result)
	if err != nil {
		t.Errorf("Failed to parse generated XML: %v", err)
	}

	if result.Name != testData.Name {
		t.Errorf("Expected name '%s', got '%s'", testData.Name, result.Name)
	}
	if result.Value != testData.Value {
		t.Errorf("Expected value %d, got %d", testData.Value, result.Value)
	}
}

func TestXMLParser_FileOperations(t *testing.T) {
	parser := NewXMLParser(false)
	testFile := "test.xml"

	// Limpiar después de la prueba
	defer os.Remove(testFile)

	testData := TestStruct{
		Name:  "Test",
		Value: 123,
	}

	// Guardar en archivo
	err := parser.SaveToFile(testData, testFile)
	if err != nil {
		t.Errorf("SaveToFile failed: %v", err)
	}

	// Cargar desde archivo
	var result TestStruct
	err = parser.LoadFromFile(testFile, &result)
	if err != nil {
		t.Errorf("LoadFromFile failed: %v", err)
	}

	if result.Name != testData.Name {
		t.Errorf("Expected name '%s', got '%s'", testData.Name, result.Name)
	}
	if result.Value != testData.Value {
		t.Errorf("Expected value %d, got %d", testData.Value, result.Value)
	}
}
