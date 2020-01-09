package simple

import (
	"encoding/json"
)

func FormatJson(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}

func ParseJson(str string, t interface{}) error {
	return json.Unmarshal([]byte(str), t)
}
