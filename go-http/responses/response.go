package responses

import (
	"avito/utils"
	"encoding/json"
	"net/http"
)

func Status(response http.ResponseWriter, statusCode int) {
	response.WriteHeader(statusCode)
	utils.LogResponse(statusCode)
}

func JSON(response http.ResponseWriter, statusCode int, data interface{}) {
	response.WriteHeader(statusCode)
	utils.LogResponse(statusCode)
	_ = json.NewEncoder(response).Encode(data)
}

func Error(response http.ResponseWriter, statusCode int, err error) {
	response.WriteHeader(statusCode)
	utils.LogResponseError(err, statusCode)
}
