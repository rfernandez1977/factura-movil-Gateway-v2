package utils

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestXML struct {
	XMLName xml.Name `xml:"test"`
	Value   string   `xml:"value,attr"`
}

func TestSimpleXML(t *testing.T) {
	test := TestXML{
		Value: "test",
	}

	data, err := xml.MarshalIndent(test, "", "  ")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `value="test"`)
}
