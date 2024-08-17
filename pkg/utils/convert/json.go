package convert

import "encoding/json"

// ToJSON converts a model to a JSON string.
func ToJSON(model interface{}) (string, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
