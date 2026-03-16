package utils_test

import (
	"os"
	"testing"

	"github.com/Deepsayan-Das/chatter_GO/internal/utils"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Exit(m.Run())
}

// ── JWT ───────────────────────────────────────────────────────────────────────

func TestGenerateJWT_ReturnsNonEmptyToken(t *testing.T) {
	token, err := utils.GenerateJWT(42)
	if err != nil {
		t.Fatalf("GenerateJWT returned unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected a non-empty token string")
	}
}

func TestValidateToken_RoundTrip(t *testing.T) {
	const userID = 99

	token, err := utils.GenerateJWT(userID)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	got, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned unexpected error: %v", err)
	}
	if got != userID {
		t.Errorf("ValidateToken: got userID %d, want %d", got, userID)
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := utils.ValidateToken("this.is.not.a.valid.token")
	if err == nil {
		t.Fatal("expected an error for an invalid token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Generate a token with the correct secret
	token, err := utils.GenerateJWT(1)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	// Swap the secret and try to validate
	os.Setenv("JWT_SECRET", "wrong-secret")
	defer os.Setenv("JWT_SECRET", "test-secret-key")

	_, err = utils.ValidateToken(token)
	if err == nil {
		t.Fatal("expected an error when validating with wrong secret, got nil")
	}
}

// ── Password ──────────────────────────────────────────────────────────────────

func TestHashPassword_ProducesHash(t *testing.T) {
	hash, err := utils.HashPassword("mysecretpassword")
	if err != nil {
		t.Fatalf("HashPassword returned unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected a non-empty hash")
	}
	// The hash should differ from the plaintext
	if hash == "mysecretpassword" {
		t.Fatal("hash must not equal the original password")
	}
}

func TestHashPassword_IsNonDeterministic(t *testing.T) {
	// bcrypt includes a random salt so two hashes of the same password differ
	hash1, _ := utils.HashPassword("password")
	hash2, _ := utils.HashPassword("password")
	if hash1 == hash2 {
		t.Error("two hashes of the same password should not be identical (bcrypt uses random salt)")
	}
}

func TestComparePassword_CorrectPassword(t *testing.T) {
	hash, err := utils.HashPassword("correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !utils.ComparePassword(hash, "correct-horse-battery-staple") {
		t.Error("ComparePassword returned false for matching password")
	}
}

func TestComparePassword_WrongPassword(t *testing.T) {
	hash, _ := utils.HashPassword("correct")
	if utils.ComparePassword(hash, "incorrect") {
		t.Error("ComparePassword returned true for non-matching password")
	}
}
