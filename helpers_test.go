package main

import (
	"os"
	"testing"

	parser "mujsann.com/ethparser/pkg"
)

func TestIsValidAddress(t *testing.T) {
	validAddress := os.Getenv("TEST_ADDRESS")
	invalidAddress := "0xInvalidAddress"

	// test a valid address
	valid, err := parser.IsValidAddress(validAddress)
	if !valid || err != nil {
		t.Errorf("expected valid address, got valid: %v, error: %v", valid, err)
	}

	// test an invalid address
	valid, err = parser.IsValidAddress(invalidAddress)
	if valid || err == nil {
		t.Errorf("expected invalid address, got valid: %v, error: %v", valid, err)
	}
}
