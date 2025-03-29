package utils

import (
	"encoding/json"
	"go-learn/model"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, message string, success bool, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Message: message,
		Success: success,
		Data:    data,
	})

}

func SetHeaders(w http.ResponseWriter, method string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", method)
}

func ValidateMethod(requestMethod string, ogMethod string, w http.ResponseWriter) {

	if requestMethod != ogMethod {
		SendJSON(w, http.StatusMethodNotAllowed, "Invalid request method", false, nil)
		return
	}
}
