package main

import (
	"testing"
)

func TestValidateKey(t *testing.T) {
	if err := validateKey(""); err == nil {
		t.Errorf("Empty key")
	}

	if err := validateKey(":"); err != nil {
		t.Errorf("Only valid key that starts with a colon")
	}

	if err := validateKey(":test"); err == nil {
		t.Errorf("Starts with colon")
	}
}

func TestValidateEmail(t *testing.T) {
	if err := validateEmail("ab"); err == nil {
		t.Errorf("Too short")
	}

	if err := validateEmail("abc"); err == nil {
		t.Errorf("No at sign")
	}

	if err := validateEmail("a@b"); err != nil {
		t.Errorf("Shortest email")
	}
}

func TestValidatePassword(t *testing.T) {
	if err := validatePassword("abcde"); err == nil {
		t.Errorf("Too short")
	}

	if err := validatePassword("abcdef"); err != nil {
		t.Errorf("Shortest password")
	}
}
