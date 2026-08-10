package auth

import "testing"

func TestGenerateSessionToken(t *testing.T) {
	first, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken returned error: %v", err)
	}

	second, err := GenerateSessionToken()
	if err != nil {
		t.Fatalf("GenerateSessionToken returned error: %v", err)
	}

	if first == "" {
		t.Fatal("session token must not be empty")
	}

	if first == second {
		t.Fatal("session tokens should be unique")
	}
}

func TestHashSessionTokenDeterministic(t *testing.T) {
	first := HashSessionToken("token")
	second := HashSessionToken("token")

	if first != second {
		t.Fatal("same token should produce same hash")
	}

	if first == "token" {
		t.Fatal("stored token hash must not equal raw token")
	}
}
