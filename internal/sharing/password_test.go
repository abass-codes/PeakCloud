package sharing

import "testing"

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("PeakCloud-Share-Password")
	if err != nil {
		t.Fatal(err)
	}

	if hash == "PeakCloud-Share-Password" {
		t.Fatal("password must not be stored in plaintext")
	}

	if !VerifyPassword(hash, "PeakCloud-Share-Password") {
		t.Fatal("correct password should verify")
	}

	if VerifyPassword(hash, "incorrect") {
		t.Fatal("incorrect password must fail")
	}
}

func TestEmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	if err != nil {
		t.Fatal(err)
	}

	if hash != "" {
		t.Fatal("empty password should produce empty hash")
	}
}
