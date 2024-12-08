package chekers

import (
	"fmt"
	"regexp"
)

func CheckPasswordValidation(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) || !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain both upper and lower case letters")
	}
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !regexp.MustCompile(`[^a-zA-Z\d]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}
