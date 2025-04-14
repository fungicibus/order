package middleware

import "strings"

func isExcludedFromMonitoring(urlPath string) bool {
	return (urlPath == "/metrics" ||
		urlPath == "/healthcheck" ||
		strings.Contains(urlPath, "swagger"))
}
