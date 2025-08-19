package typs

import (
	// "math/rand"
	"crypto/rand"
	`math/big`
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

// The following constants are to be used only by this package and are thus not exported

const smallLetterBytes = "abcdefghijklmnopqrstuvwxyz"
const capitalLetterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digitBytes = "1234567890"
const specialCharBytes = "/}%,]'?!<_~@;{`&:^+|.*#$()-=[>\""

// GetRandomAlphaString will get a n character long random alphabetic string
func GetRandomAlphaString(n int) string {
	// SECURE: uses crypto/rand
	letterBytes := smallLetterBytes + capitalLetterBytes
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[num.Int64()]
	}
	return string(b)
}

// GetRandomNumericString will get a n character long random numeric string
func GetRandomNumericString(n int) string {
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digitBytes))))
		b[i] = digitBytes[num.Int64()]
	}
	return string(b)
}

// GetRandomSpecialCharacterString will get a n character long random string made of special characters
func GetRandomSpecialCharacterString(n int) string {
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(specialCharBytes))))
		b[i] = specialCharBytes[num.Int64()]
	}
	return string(b)
}

// GetRandomString will get a n character long random string that might contain alphabets, numbers and special chars
func GetRandomString(n int) string {
	s := capitalLetterBytes + smallLetterBytes + digitBytes + specialCharBytes
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(s))))
		b[i] = s[num.Int64()]
	}
	return string(b)
}

func CountCapitalLetters(str string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		if strings.ContainsAny(string(str[i]), capitalLetterBytes) {
			count += 1
		}
	}
	return count
}

func CountSmallLetters(str string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		if strings.ContainsAny(string(str[i]), smallLetterBytes) {
			count += 1
		}
	}
	return count
}

func CountNumericCharacters(str string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		if strings.ContainsAny(string(str[i]), digitBytes) {
			count += 1
		}
	}
	return count
}

func CountSpecialCharacters(str string) int {
	count := 0
	for i := 0; i < len(str); i++ {
		if strings.ContainsAny(string(str[i]), specialCharBytes) {
			count += 1
		}
	}
	return count
}

func IsAlphaNumeric(input string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return pattern.MatchString(input)
}

func IsAlphaNumericOrDotDashUnderscore(input string) bool {
	pattern := regexp.MustCompile("^[a-zA-Z0-9\\._-]+$")
	return pattern.MatchString(input)
}

func IsPositiveNumber(input string) bool {
	pattern := regexp.MustCompile("^[0-9]+$")
	return pattern.MatchString(input)
}

// IsValidEmail checks if the supplied string is a valid email
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	// The standard library will allow an email like `abcd@example` to pass because `example` can be
	// a valid hostname defined on the localhost (/etc/hosts) or local network.
	// For a fully valid email address, we must check that the input ends in a TLD

	validTLD := regexp.MustCompile(`\.[a-z]{2,}$`)
	return validTLD.MatchString(email)
}

// Marker: =======================================
// Marker: String type
// Marker: =======================================

// String facilitates basic string type and functions
type String struct {
	originalValue string // The value with which this type was created
	currentValue  string // The value with which the functions will work
}

// NewString returns a new String type
func NewString(str string) *String {
	return &String{
		originalValue: str,
		currentValue:  str,
	}
}

// String implements the built-in Stringer interface
func (s *String) String() string {
	return s.currentValue
}

// Reset sets the current value of the String to its original value with which it was created
func (s *String) Reset() {
	s.currentValue = s.originalValue
}

// GetOriginalString returns the original value with which the String wqs created
func (s String) GetOriginalString() string {
	return s.originalValue
}

// Length returns the length of the string
func (s *String) Length() int {
	return len(s.currentValue)
}

// Trim trims off spaces from both ends of the string
func (s *String) Trim() *String {
	s.currentValue = strings.TrimSpace(s.currentValue)
	return s
}

// ToUpperCase converts the string to lower case
func (s *String) ToUpperCase() *String {
	s.currentValue = strings.ToUpper(s.currentValue)
	return s
}

// ToLowerCase converts the string to lower case
func (s *String) ToLowerCase() *String {
	s.currentValue = strings.ToLower(s.currentValue)
	return s
}

// IsNumeric tells if a string is numeric or not
func (s *String) IsNumeric() bool {
	if _, err := strconv.Atoi(s.currentValue); err == nil {
		return true
	}

	return false
}
