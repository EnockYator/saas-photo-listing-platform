package middleware

import "strings"

func Sanitize(s string) string {
	// remove CRLF and other control chars that enable log injection
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}