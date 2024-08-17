package locationManager

import (
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type Sanitizer struct {
	value string
}

func NewSanitizer(s string) *Sanitizer {
	sanitizePolicy := bluemonday.UGCPolicy()
	trimmedString := strings.TrimSpace(s)
	trimmedAndSanitizedString := sanitizePolicy.Sanitize(trimmedString)
	return &Sanitizer{value: trimmedAndSanitizedString}
}

func (s *Sanitizer) MaxLength(maxLength int) *Sanitizer {
	if len(s.value) > maxLength {
		s.value = s.value[:maxLength]
	}
	return s
}

func (s *Sanitizer) ValidateURL() *Sanitizer {
	if _, err := url.ParseRequestURI(s.value); err != nil {
		s.value = ""
	}
	return s
}

func (s *Sanitizer) Result() string {
	return s.value
}
