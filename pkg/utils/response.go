package utils

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

func ResponseJson(writer http.ResponseWriter, responseData interface{}, statusCode int) error {
	data, err := json.Marshal(responseData)

	if err != nil {
		return err
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	_, err = writer.Write(data)

	return err
}

func ResponseCode(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
}
