package utils

import "golang.org/x/crypto/bcrypt"

func CreateHashFromText(text string, cost int32) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(text), int(cost))

	hashedPassword := string(hash)

	return hashedPassword, err
}
