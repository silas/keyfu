package main

import (
	"errors"
	"strings"
)

func validateKey(key string) (err error) {
	if key == "" {
		err = errors.New("Non-empty key required.")
	} else if key != ":" && key[0] == ':' {
		err = errors.New("Keys starting with a colon are reserved.")
	}

	return err
}

func validateEmail(email string) (err error) {
	if len(email) < 3 || strings.Index(email, "@") < 1 {
		err = errors.New("Invalid email address.")
	}

	return err
}

func validatePassword(password string) (err error) {
	if len(password) < 6 {
		err = errors.New("Password must be at least 6 characters.")
	}

	return err
}
