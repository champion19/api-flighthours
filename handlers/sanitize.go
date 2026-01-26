package handlers

import "strings"

func TrimString(s string) string {
	return strings.TrimSpace(s)
}

func TrimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	return &trimmed
}

type Sanitizable interface {
	Sanitize()
}

