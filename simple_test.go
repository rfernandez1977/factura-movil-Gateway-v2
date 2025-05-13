package main

import (
	"fmt"
	"testing"
)

// TestSimple es una prueba simple para verificar que Go funciona
func TestSimple(t *testing.T) {
	// Esta es una prueba simple para verificar que Go funciona
	expected := 4
	actual := 2 + 2

	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	} else {
		fmt.Println("Simple test passed!")
	}
}
