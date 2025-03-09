package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	plaintext := "password"

	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Error("expected hash to not be empty")
	}

	if hash == plaintext {
		t.Error("expected hash to be different from plaintext")
	}
}

func TestComparePassword(t *testing.T) {
	plaintext := "password"
	fail_plaintext := "not password"

	hash, err := HashPassword(plaintext)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !ComparePassword(hash, []byte(plaintext)) {
		t.Errorf("error password to match hash")
	}

	if ComparePassword(hash, []byte(fail_plaintext)) {
		t.Errorf("error password to not match hash")
	}

}
