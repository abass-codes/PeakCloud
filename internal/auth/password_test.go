package auth

import "testing"

func TestPasswordHashAndVerify(t *testing.T) {
	password := "PeakCloud-Test-Password-123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatal("password hash must not equal plaintext password")
	}

	valid, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned error: %v", err)
	}

	if !valid {
		t.Fatal("expected password to verify")
	}
}

func TestPasswordRejectsIncorrectPassword(t *testing.T) {
	hash, err := HashPassword("Correct-Password-123!")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	valid, err := VerifyPassword("Wrong-Password-123!", hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned error: %v", err)
	}

	if valid {
		t.Fatal("incorrect password must not verify")
	}
}
