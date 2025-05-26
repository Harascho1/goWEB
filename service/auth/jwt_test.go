package auth

import (
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Errorf("Expected token to be not empty")
	}
}
