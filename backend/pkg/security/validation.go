package security

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrRequired     = errors.New("field is required")
	ErrTooShort     = errors.New("value too short")
	ErrTooLong      = errors.New("value too long")
	ErrInvalidEmail = errors.New("invalid email address")
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrInvalidURL   = errors.New("invalid URL")
	ErrWeakPassword = errors.New("password too weak")
	ErrInvalidInput = errors.New("invalid input")
)

// ValidationResult holds validation results
type ValidationResult struct {
	Valid  bool
	Errors map[string]string
}

// NewValidationResult creates an empty result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: make(map[string]string),
	}
}

// AddError adds a validation error
func (v *ValidationResult) AddError(field, message string) {
	v.Valid = false
	v.Errors[field] = message
}

// Validator provides input validation
type Validator struct {
	result *ValidationResult
}

// NewValidator creates a validator
func NewValidator() *Validator {
	return &Validator{
		result: NewValidationResult(),
	}
}

// Result returns the validation result
func (v *Validator) Result() *ValidationResult {
	return v.result
}

// Required validates field is not empty
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.result.AddError(field, "is required")
	}
	return v
}

// MinLength validates minimum length
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.result.AddError(field, "must be at least "+string(rune(min+'0'))+" characters")
	}
	return v
}

// MaxLength validates maximum length
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.result.AddError(field, "must be at most "+string(rune(max+'0'))+" characters")
	}
	return v
}

// Email validates email format
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.result.AddError(field, "invalid email format")
	}
	return v
}

// Phone validates phone number
func (v *Validator) Phone(field, value string) *Validator {
	if value == "" {
		return v
	}
	// Vietnamese phone number pattern
	pattern := regexp.MustCompile(`^(\+84|0)(3|5|7|8|9)[0-9]{8}$`)
	if !pattern.MatchString(value) {
		v.result.AddError(field, "invalid phone number")
	}
	return v
}

// Password validates password strength
func (v *Validator) Password(field, value string) *Validator {
	if len(value) < 8 {
		v.result.AddError(field, "password must be at least 8 characters")
		return v
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range value {
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
		v.result.AddError(field, "password must contain uppercase, lowercase, number, and special character")
	}
	return v
}

// Alphanumeric validates alphanumeric only
func (v *Validator) Alphanumeric(field, value string) *Validator {
	if value == "" {
		return v
	}
	pattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !pattern.MatchString(value) {
		v.result.AddError(field, "must contain only letters and numbers")
	}
	return v
}

// NoSQLInjection checks for SQL injection patterns
func (v *Validator) NoSQLInjection(field, value string) *Validator {
	dangerous := []string{
		"--", ";--", "/*", "*/", "@@", "@",
		"'", "\"", "`",
		"UNION", "SELECT", "INSERT", "UPDATE", "DELETE", "DROP",
		"EXEC", "EXECUTE", "xp_", "sp_",
		"OR 1=1", "OR '1'='1", "' OR ''='",
	}

	upper := strings.ToUpper(value)
	for _, pattern := range dangerous {
		if strings.Contains(upper, strings.ToUpper(pattern)) {
			v.result.AddError(field, "contains invalid characters")
			return v
		}
	}
	return v
}

// NoXSS checks for XSS patterns
func (v *Validator) NoXSS(field, value string) *Validator {
	dangerous := []string{
		"<script", "</script>", "javascript:",
		"onerror=", "onload=", "onclick=", "onmouseover=",
		"<iframe", "<object", "<embed", "<svg",
	}

	lower := strings.ToLower(value)
	for _, pattern := range dangerous {
		if strings.Contains(lower, pattern) {
			v.result.AddError(field, "contains invalid content")
			return v
		}
	}
	return v
}

// SanitizeInput removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Replace multiple spaces with single space
	space := regexp.MustCompile(`\s+`)
	input = space.ReplaceAllString(input, " ")
	
	return input
}

// SanitizeHTML escapes HTML special characters
func SanitizeHTML(input string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(input)
}

// SanitizeSQL escapes SQL special characters (use parameterized queries instead!)
func SanitizeSQL(input string) string {
	replacer := strings.NewReplacer(
		"'", "''",
		"\\", "\\\\",
	)
	return replacer.Replace(input)
}

// ValidateSlug validates URL slug format
func ValidateSlug(slug string) bool {
	pattern := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	return pattern.MatchString(slug)
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) bool {
	pattern := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return pattern.MatchString(uuid)
}
