package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash checks if a password matches its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// PasswordComplexityError represents password validation errors
type PasswordComplexityError struct {
	Message string
	Errors  []string
}

func (e *PasswordComplexityError) Error() string {
	return e.Message
}

// PasswordComplexityRequirements defines password complexity rules
type PasswordComplexityRequirements struct {
	MinLength        int
	MaxLength        int
	RequireUppercase bool
	RequireLowercase bool
	RequireNumbers   bool
	RequireSpecial   bool
	ForbiddenWords   []string
}

// DefaultPasswordRequirements returns default password complexity requirements
func DefaultPasswordRequirements() PasswordComplexityRequirements {
	return PasswordComplexityRequirements{
		MinLength:        8,
		MaxLength:        128,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   true,
		ForbiddenWords:   []string{"password", "123456", "admin", "user", "login"},
	}
}

// ValidatePasswordComplexity validates password against complexity requirements
func ValidatePasswordComplexity(password string, requirements ...PasswordComplexityRequirements) error {
	req := DefaultPasswordRequirements()
	if len(requirements) > 0 {
		req = requirements[0]
	}

	var errors []string

	if len(password) < req.MinLength {
		errors = append(errors, fmt.Sprintf("Password must be at least %d characters long", req.MinLength))
	}

	if len(password) > req.MaxLength {
		errors = append(errors, fmt.Sprintf("Password must not exceed %d characters", req.MaxLength))
	}

	if req.RequireUppercase {
		hasUpper := false
		for _, char := range password {
			if unicode.IsUpper(char) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			errors = append(errors, "Password must contain at least one uppercase letter")
		}
	}

	if req.RequireLowercase {
		hasLower := false
		for _, char := range password {
			if unicode.IsLower(char) {
				hasLower = true
				break
			}
		}
		if !hasLower {
			errors = append(errors, "Password must contain at least one lowercase letter")
		}
	}

	if req.RequireNumbers {
		hasNumber := false
		for _, char := range password {
			if unicode.IsDigit(char) {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			errors = append(errors, "Password must contain at least one number")
		}
	}

	if req.RequireSpecial {
		specialChars := `!@#$%^&*()_+-=[]{}|;:,.<>?`
		hasSpecial := false
		for _, char := range password {
			for _, special := range specialChars {
				if char == special {
					hasSpecial = true
					break
				}
			}
			if hasSpecial {
				break
			}
		}
		if !hasSpecial {
			errors = append(errors, "Password must contain at least one special character (!@#$%^&*()_+-=[]{}|;:,.<>?)")
		}
	}

	passwordLower := strings.ToLower(password)
	for _, forbidden := range req.ForbiddenWords {
		if strings.Contains(passwordLower, strings.ToLower(forbidden)) {
			errors = append(errors, fmt.Sprintf("Password must not contain common words like '%s'", forbidden))
		}
	}

	if matched, _ := regexp.MatchString(`(.)\1{2,}`, password); matched {
		errors = append(errors, "Password must not contain repeated characters (e.g., 'aaa', '111')")
	}

	if matched, _ := regexp.MatchString(`(012|123|234|345|456|567|678|789|890|abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`, strings.ToLower(password)); matched {
		errors = append(errors, "Password must not contain sequential characters (e.g., '123', 'abc')")
	}

	if len(errors) > 0 {
		return &PasswordComplexityError{
			Message: "Password does not meet complexity requirements",
			Errors:  errors,
		}
	}

	return nil
}

// ValidateAndHashPassword validates password complexity and returns hashed password
func ValidateAndHashPassword(password string, requirements ...PasswordComplexityRequirements) (string, error) {
	if err := ValidatePasswordComplexity(password, requirements...); err != nil {
		return "", err
	}
	return HashPassword(password)
}
