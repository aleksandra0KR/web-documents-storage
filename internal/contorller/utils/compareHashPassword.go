package utils

import "golang.org/x/crypto/bcrypt"

func CompareHashPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
