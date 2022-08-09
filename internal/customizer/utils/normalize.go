package utils

import (
	"alnovi/customizer/pkg/flatten"
	"strings"
)

func PreparedStoreData(data *flatten.Flatten) *flatten.Flatten {
	newData := flatten.NewFlatten()
	newData.SetDelimiter(".")

	for key, val := range data.All(true) {
		subKeys := strings.Split(key, ".")
		sabKeysVal := subKeys[len(subKeys)-1:][0]

		switch sabKeysVal {
		case "description":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+data.GetDelimiter()+sabKeysVal, val)
		case "value":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+data.GetDelimiter()+sabKeysVal, val)
		default:
			newData.Add(key, map[string]interface{}{
				"description": "",
				"value":       val,
			})
		}
	}

	return newData
}

func PreparedViewData(data *flatten.Flatten) *flatten.Flatten {
	newData := flatten.NewFlatten()

	for key, val := range data.All(true) {
		subKeys := strings.Split(key, "#")
		sabKeysVal := subKeys[len(subKeys)-1:][0]

		switch sabKeysVal {
		case "value":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], "."), val)
		default:
			break
		}
	}

	return newData
}

func DataMongoNormalize(data *flatten.Flatten) *flatten.Flatten {
	newData := flatten.NewFlatten()
	newData.SetDelimiter("#")

	for key, val := range data.All(true) {
		subKeys := strings.Split(key, ".")
		sabKeysVal := subKeys[len(subKeys)-1:][0]

		switch sabKeysVal {
		case "description":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+"#"+sabKeysVal, val)
		case "value":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+"#"+sabKeysVal, val)
		default:
			newData.Add(key, map[string]interface{}{
				"description": "",
				"value":       val,
			})
		}
	}

	return newData
}

func DataMongoDenormalize(data *flatten.Flatten) *flatten.Flatten {
	newData := flatten.NewFlatten()
	newData.SetDelimiter(".")

	for key, val := range data.All(true) {
		subKeys := strings.Split(key, "#")
		sabKeysVal := subKeys[len(subKeys)-1:][0]

		switch sabKeysVal {
		case "description":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+"."+sabKeysVal, val)
		case "value":
			newData.Add(strings.Join(subKeys[:len(subKeys)-1], ".")+"."+sabKeysVal, val)
		default:
			newData.Add(key, map[string]interface{}{
				"description": "",
				"value":       val,
			})
		}
	}

	return newData
}

func ClearArtifactFromStruct(data *flatten.Flatten) *flatten.Flatten {
	for key, _ := range data.All(true) {
		subKeys := strings.Split(key, ".")
		sabKeysNamespace := strings.Join(subKeys[:len(subKeys)-1], ".")

		if len(data.Keys(sabKeysNamespace)) == 2 {

			continue
		}

		data.Delete(key)
	}

	return data
}
