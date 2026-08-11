package sharing

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(hash, password string) bool {
	if hash == "" {
		return true
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	) == nil
}
