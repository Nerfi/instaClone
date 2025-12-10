package security

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// utils functions for password security and hashing

// En otras palabras: un método con receiver = un método de una clase.
const (
	DefaultMinLength = 8
	DefaultMaxLength = 72
)

func HashPassword(password string) (string, error) {
	return HashPasswordWithLimits(password, DefaultMinLength, DefaultMaxLength)
}

func HashPasswordWithLimits(password string, minLength, maxLenght int) (string, error) {
	// validate the password
	if err := ValidatePassword(password, minLength, maxLenght); err != nil {
		return "", err
	}

	// generate hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ValidatePassword(password string, minLength, maxLength int) error {
	if len(password) < minLength {
		return fmt.Errorf("password length should be at least %d characters", minLength)
	}

	if len(password) > maxLength {
		return fmt.Errorf("password length should be at most %d characters", maxLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}

func ComparePassword(hash, password string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil

}
