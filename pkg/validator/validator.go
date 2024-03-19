package validator

import "fmt"

func ValidateString(value string, minLength, maxlength int) error {
	n := len(value)
	if n < minLength || n > maxlength {
		return fmt.Errorf("must contain from %d-%d", minLength, maxlength)
	}
	return nil
}

func ValidateTitle(value string) error {
	if err := ValidateString(value, 3, 255); err != nil {
		return err
	}
	return nil
}
