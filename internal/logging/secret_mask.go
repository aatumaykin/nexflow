package logging

import "strings"

// secretFields contains patterns that identify secret fields
var secretFields = map[string]bool{
	"api_key":       true,
	"apikey":        true,
	"apiKey":        true,
	"token":         true,
	"access_token":  true,
	"accesstoken":   true,
	"accessToken":   true,
	"refresh_token": true,
	"refreshtoken":  true,
	"refreshToken":  true,
	"password":      true,
	"pass":          true,
	"secret":        true,
	"private_key":   true,
	"privatekey":    true,
	"privateKey":    true,
	"auth_token":    true,
	"authtoken":     true,
	"authToken":     true,
	"bot_token":     true,
	"bottoken":      true,
	"botToken":      true,
	"bearer_token":  true,
	"bearertoken":   true,
	"bearerToken":   true,
	"client_secret": true,
	"clientsecret":  true,
	"clientSecret":  true,
	"api_secret":    true,
	"apisecret":     true,
	"apiSecret":     true,
	"session_token": true,
	"sessiontoken":  true,
	"sessionToken":  true,
	"csrf_token":    true,
	"csrftoken":     true,
	"csrfToken":     true,
	"authorization": true,
	"credentials":   true,
	"credential":    true,
}

// shouldMask checks if a field key should have its value masked
func shouldMask(key string) bool {
	lowerKey := strings.ToLower(key)
	return secretFields[lowerKey]
}

// maskValue masks a secret value by showing only first and last characters
func maskValue(value string) string {
	if len(value) <= 4 {
		return "***"
	}
	// Show first 2 and last 2 characters, mask the rest
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}
