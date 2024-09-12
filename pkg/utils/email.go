package utils

import "regexp"

// ValidateEmailFormat validates the format of an email address using a regular expression
func ValidateEmailFormat(email string) bool {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailRegexPattern)

	return emailRegex.MatchString(email)
}
