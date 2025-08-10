package utils

import (
	"reflect"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
	nonstd_validators "github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/google/uuid"
	"github.com/nbrglm/auth-platform/opts"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", ValidatePasswordRequirements)
	v.RegisterValidation("notblank", nonstd_validators.NotBlank)
	v.RegisterValidation("uuidv7", ValidateUUIDV7)
	v.RegisterValidation("domain", ValidateDomain)
}

var Validator *validator.Validate

func InitValidator() {
	if Validator != nil {
		return
	}
	Validator = validator.New()
	RegisterCustomValidators(Validator)
}

// Validates whether the field is a valid UUID v7
func ValidateUUIDV7(fl validator.FieldLevel) bool {
	val := ""
	field := fl.Field()
	switch field.Kind() {
	case reflect.Ptr:
		valType := field.Elem().Kind()
		if valType != reflect.String {
			return false
		}
		val = field.Elem().String()
	case reflect.String:
		val = field.String()
	default:
		return false
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return false
	}

	return id.Version() == 7
}

var allowedSpecialCharactersForPassword = "-_*@."

func ValidatePasswordRequirements(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	length := len(password)

	if length < 8 || length > 32 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecialSymbol bool

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case containsRunes(allowedSpecialCharactersForPassword, c):
			hasSpecialSymbol = true
		case unicode.IsSpace(c):
			// Do not allow spaces in passwords
			return false
		default:
			// We explicitly disallow other characters
			// This helps in preventing any kind of injection attacks
			// as well as increases the performance
			return false
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecialSymbol
}

func containsRunes(runes string, runeToFind rune) bool {
	for _, r := range runes {
		if r == runeToFind {
			return true
		}
	}
	return false
}

// domainRegex is a regex pattern to validate domain names.
var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,63}$`)

// ValidateDomain checks if the provided field is a valid domain.
//
// Only valid format is: "example.com" or "sub.example.com"
func ValidateDomain(fl validator.FieldLevel) bool {
	domain := fl.Field().String()
	if domain == "" {
		return false
	}

	if opts.Debug && domain == "localhost" {
		return true
	}

	// A simple check for domain format
	if len(domain) < 3 || len(domain) > 253 {
		return false
	}

	return domainRegex.MatchString(domain)
}
