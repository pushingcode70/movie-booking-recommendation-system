package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain text password into a bcrypt hash.
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

// CheckPassword compares a plain password with its bcrypt hash.
func CheckPassword(password string, hashedPassword string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

//extracts plain text from user entered and extracts salt from db
// and we add them it will give what stored as hashed password in db
//It takes the plain password the user entered.
//It extracts the salt (and cost) from the stored bcrypt hash.
// It runs the bcrypt algorithm again using that password and the extracted salt.
// If the newly computed hash matches the stored hash, login succeeds.

//like we had email xyz gmail.com and it has hashed password stored and user enters
//plain text it will extract that text and add salt taken form hashedpassword already
//stored in db so like abc123 + salt then it takes as whole and compares with db which
//has the required email ie abc123
//remember it when password is entered than all this happens then gets compared wih db hashed password
