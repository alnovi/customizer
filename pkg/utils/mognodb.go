package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NormalizeMongoPrimitive(data interface{}) interface{} {
	var result interface{}

	switch data.(type) {
	case primitive.D:
		subResult := map[string]interface{}{}

		for key, val := range data.(primitive.D).Map() {
			subResult[key] = NormalizeMongoPrimitive(val)
		}

		result = subResult

	case primitive.M:
		subResult := map[string]interface{}{}

		for key, val := range data.(primitive.M) {
			subResult[key] = NormalizeMongoPrimitive(val)
		}

		result = subResult

	case primitive.A:
		subResult := make([]interface{}, len(data.(primitive.A)))

		for key, val := range data.(primitive.A) {
			subResult[key] = NormalizeMongoPrimitive(val)
		}

		result = subResult
	default:
		result = data
	}

	return result
}
