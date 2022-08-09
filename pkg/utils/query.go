package utils

import "net/http"

func QueryValue(request *http.Request, field string) V {
	query := request.URL.Query()

	if val, ok := query[field]; ok && val[0] != "" {
		return NewV(val[0])
	}

	return NewV(EmptyValue{})
}
