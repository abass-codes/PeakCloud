package sharing

import "testing"

func TestGenerateToken(t *testing.T) {
	first, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	second, err := GenerateToken()
	if err != nil {
		t.Fatal(err)
	}

	if first == "" || second == "" {
		t.Fatal("tokens cannot be empty")
	}

	if first == second {
		t.Fatal("tokens must be unique")
	}
}

func TestHashTokenDoesNotExposeToken(t *testing.T) {
	token := "secret-token"
	hash := HashToken(token)

	if hash == token {
		t.Fatal("token must not be stored directly")
	}
}

func TestHashTokenDeterministic(t *testing.T) {
	if HashToken("abc") != HashToken("abc") {
		t.Fatal("hash must be deterministic")
	}
}
