package chekers

import (
	"fmt"
	"regexp"
)

func CheckLoginValidation(login string) error {
	if len(login) < 8 || !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(login) {
		return fmt.Errorf("invalid login")
	}
	return nil
}
