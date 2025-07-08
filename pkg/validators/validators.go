package validators

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Custom Validation Func
func RegisterCustomValidations() {
	_ = validate.RegisterValidation("role", validateRole)
}

func validateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	return role == "super_admin" || role == "admin" || role == "user"
}
