package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashedPassword(password string) (h string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil

}

func VerifyPassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
