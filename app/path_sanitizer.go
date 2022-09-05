package app

import (
	"net/url"
	"path/filepath"
	"strings"
)

//
// SanitizePath
// @Description: Sanitize a webpack url / path
// @param str string
// @return string
func SanitizePath(str string) string {
	if u, err := url.Parse(str); err == nil && u.Path != "" {
		str = u.Path
	}
	if strings.Contains(str, " ") {
		str = strings.Split(str, " ")[0]
	}

	str = filepath.Clean(str)
	result := make([]string, 0)
	for _, p := range strings.Split(str, "/") {
		if p == ".." {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else {
			result = append(result, p)
		}
	}

	return strings.Join(result, "/")
}
