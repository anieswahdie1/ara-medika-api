package validators

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New()
	registerCustomValidations()
}

// Custom Validation Func
func registerCustomValidations() {
	// Validasi custom untuk role
	_ = Validate.RegisterValidation("role", validateRole)
	_ = Validate.RegisterValidation("strong_password", validateStrongPassword)
}

func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == "super_admin" || role == "admin" || role == "user"
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimal 8 karakter, mengandung angka, huruf besar dan kecil
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
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

	return hasUpper && hasLower && hasNumber && hasSpecial
}
