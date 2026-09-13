package utils

import "golang.org/x/crypto/bcrypt"

// hashPassword converts a plain text password into a bcrypt hash
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPassword compares a plain password with its bcrypt hash
func CheckPassword(password string, hashedPassword string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
