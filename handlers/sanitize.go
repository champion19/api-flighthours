package handlers

import "strings"

func TrimString(s string) string {
	return strings.TrimSpace(s)
}

//TrimStringPtr trims a string pointer,returns nil if the pointer is nil
func TrimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

//sanitizable interface for structs that can be sanitized
type Sanitizable interface {
	Sanitize()
}

