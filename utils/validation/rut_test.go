package validation

import (
	"fmt"
	"testing"
)

func TestValidateRUT(t *testing.T) {
	tests := []struct {
		name    string
		rut     string
		wantErr bool
	}{
		{
			name:    "RUT válido empresa con guión",
			rut:     "76212889-6",
			wantErr: false,
		},
		{
			name:    "RUT válido persona natural con guión",
			rut:     "13195458-1",
			wantErr: false,
		},
		{
			name:    "RUT válido con K",
			rut:     "10138666-K",
			wantErr: false,
		},
		{
			name:    "RUT válido con puntos y guión",
			rut:     "76.212.889-6",
			wantErr: false,
		},
		{
			name:    "RUT válido sin formato",
			rut:     "131954581",
			wantErr: false,
		},
		{
			name:    "RUT inválido - dígito verificador incorrecto",
			rut:     "76212889-5",
			wantErr: true,
		},
		{
			name:    "RUT inválido - formato incorrecto",
			rut:     "76212889X",
			wantErr: true,
		},
		{
			name:    "RUT vacío",
			rut:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRUT(tt.rut)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRUT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalcularDV(t *testing.T) {
	tests := []struct {
		name   string
		rut    int
		wantDV string
	}{
		{
			name:   "RUT empresa 76212889",
			rut:    76212889,
			wantDV: "6",
		},
		{
			name:   "RUT persona natural 13195458",
			rut:    13195458,
			wantDV: "1",
		},
		{
			name:   "RUT con K 10138666",
			rut:    10138666,
			wantDV: "K",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcularDV(tt.rut); got != tt.wantDV {
				t.Errorf("calcularDV() = %v, want %v", got, tt.wantDV)
			}
		})
	}
}

func TestCalcularDVPasoAPaso(t *testing.T) {
	rut := 10138666
	rutStr := "10138666"
	var suma int
	var multiplicador = 2

	// Calcular suma
	for i := len(rutStr) - 1; i >= 0; i-- {
		digito := int(rutStr[i] - '0')
		producto := digito * multiplicador
		t.Logf("Dígito %d * %d = %d", digito, multiplicador, producto)
		suma += producto
		multiplicador++
		if multiplicador > 7 {
			multiplicador = 2
		}
	}

	t.Logf("Suma total: %d", suma)
	resto := suma % 11
	t.Logf("Resto = %d %% 11 = %d", suma, resto)

	var dv string
	if resto == 0 {
		dv = "0"
	} else if resto == 1 {
		dv = "K"
	} else {
		dv = fmt.Sprintf("%d", 11-resto)
	}

	t.Logf("Dígito verificador calculado: %s", dv)
	t.Logf("Dígito verificador esperado: K")

	if got := calcularDV(rut); got != "K" {
		t.Errorf("calcularDV() = %v, want K", got)
	}
}
