package unit

import (
	"testing"

	"FMgo/pkg/sii"
)

func TestSIIAuthentication(t *testing.T) {
	tests := []struct {
		name     string
		certPath string
		keyPath  string
		wantErr  bool
	}{
		{
			name:     "valid_cert",
			certPath: "../certs/test_client.crt",
			keyPath:  "../certs/test_client.key",
			wantErr:  false,
		},
		{
			name:     "invalid_cert",
			certPath: "nonexistent.crt",
			keyPath:  "nonexistent.key",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := sii.NewClient(tt.certPath, tt.keyPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				err = client.Authenticate()
				if err != nil {
					t.Errorf("Authenticate() error = %v", err)
				}
			}
		})
	}
}

func TestSIIEnvio(t *testing.T) {
	tests := []struct {
		name     string
		dteXML   string
		wantErr  bool
		wantCode string
	}{
		{
			name: "envio_valido",
			dteXML: `<?xml version="1.0" encoding="UTF-8"?>
<DTE version="1.0">
    <!-- Contenido de prueba -->
</DTE>`,
			wantErr:  false,
			wantCode: "0",
		},
		{
			name:     "xml_invalido",
			dteXML:   "contenido no xml",
			wantErr:  true,
			wantCode: "",
		},
	}

	// Crear cliente SII para las pruebas
	client, err := sii.NewClient("../certs/test_client.crt", "../certs/test_client.key")
	if err != nil {
		t.Fatalf("error creando cliente SII: %v", err)
	}

	// Autenticar cliente
	if err := client.Authenticate(); err != nil {
		t.Fatalf("error autenticando cliente: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trackID, err := client.SendDTE(tt.dteXML)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendDTE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && trackID == "" {
				t.Error("SendDTE() trackID vacío")
			}
		})
	}
}
