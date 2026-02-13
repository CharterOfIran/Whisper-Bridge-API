package security

import "strings"

func Clean(msg string) string {
	clean := strings.ReplaceAll(msg, "<", "&lt;")
	clean = strings.ReplaceAll(clean, ">", "&gt;")
	if len(clean) > 5000 {
		return clean[:5000]
	}
	return clean
}
