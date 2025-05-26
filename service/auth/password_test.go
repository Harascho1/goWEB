package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Error(err)
	}

	if hash == "" {
		t.Errorf("expected hash to be not empty")
	}

	if hash == "password" {
		t.Errorf("expected hast to be diffrent from password")
	}
}

func TestComparePassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Error(err)
	}

	if !ComparePassword(hash, []byte("password")) {
		t.Errorf("expected password to match hash")
	}

	if !ComparePassword(hash, []byte("notpassword")) {
		t.Errorf("expected password to not match hash")
	}
}
