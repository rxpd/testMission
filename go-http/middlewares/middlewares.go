package middlewares

import (
	"avito/utils"
	"net/http"
)

func SetJSONmdw(next http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		utils.LogRequest(request)
		setHeader(response)
		if request.Method == http.MethodOptions {
			response.WriteHeader(http.StatusOK)
			return
		}
		next(response, request)
	}
}



func setHeader(response http.ResponseWriter) {
	//response.Header().Set("Access-Control-Allow-Origin", "origin")
	// 	response.Header().Set("Access-Control-Allow-Headers", "Authorization")
	response.Header().Set("Access-Control-Allow-Headers", "*")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Expose-Headers", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Content-Security-Policy", "default-src 'self'")
	response.Header().Set("Strict-Transport-Security", "max-age=15552000 [; includeSubdomains]")
	response.Header().Set("X-Frame-Options", "deny")
	response.Header().Set("X-Content-Type-Options", "nosniff")
	response.Header().Set("X-XSS-Protection", "1; mode=block")
	response.Header().Set("Referrer-Policy", "no-referrer")
}
