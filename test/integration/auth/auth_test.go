package auth_test

import (
	"context"
	"strings"
	"testing"

	"FMgo/core/sii"
)

func TestSemillaIntegration(t *testing.T) {
	// Configurar cliente en modo de prueba
	config := &sii.Config{
		BaseURL:     "https://palena.sii.cl",
		Environment: "certificacion",
		Timeout:     30,
		MaxRetries:  5,
		TestMode:    true,
	}
	client := sii.NewClient(config)
	ctx := context.Background()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Obtener semilla exitosamente",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			semilla, err := client.GetSemilla(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSemilla() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if semilla == "" {
					t.Error("GetSemilla() retornó semilla vacía")
				}
				if !strings.Contains(semilla, "SEMILLA-DE-PRUEBA") {
					t.Errorf("GetSemilla() retornó semilla inesperada: %s", semilla)
				}
				t.Logf("Semilla obtenida: %s", semilla)
			}
		})
	}
}

func TestTokenIntegration(t *testing.T) {
	// Configurar cliente en modo de prueba
	config := &sii.Config{
		BaseURL:     "https://palena.sii.cl",
		Environment: "certificacion",
		Timeout:     30,
		MaxRetries:  5,
		TestMode:    true,
	}
	client := sii.NewClient(config)
	ctx := context.Background()

	tests := []struct {
		name         string
		semilla      string
		wantErr      bool
		errorMessage string
	}{
		{
			name:    "Obtener token exitosamente",
			semilla: "SEMILLA-VALIDA",
			wantErr: false,
		},
		{
			name:         "Token con semilla inválida",
			semilla:      "SEMILLA-INVALIDA",
			wantErr:      true,
			errorMessage: "error SOAP: SOAP-ENV:Client - Semilla inválida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Primero obtenemos una semilla real si es necesario
			var semilla string
			if tt.semilla == "SEMILLA-VALIDA" {
				var err error
				semilla, err = client.GetSemilla(ctx)
				if err != nil {
					t.Fatalf("No se pudo obtener semilla: %v", err)
				}
			} else {
				semilla = tt.semilla
			}

			// Intentamos obtener el token
			token, err := client.GetToken(ctx, semilla)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if token == "" {
					t.Error("GetToken() retornó token vacío")
				}
				if !strings.Contains(token, "TOKEN-DE-PRUEBA") {
					t.Errorf("GetToken() retornó token inesperado: %s", token)
				}
				t.Logf("Token obtenido: %s", token)
			}
			if tt.wantErr && err != nil && tt.errorMessage != "" {
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("GetToken() error message = %v, want contener %v", err.Error(), tt.errorMessage)
				}
			}
		})
	}
}
