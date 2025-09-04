package jsons

import (
	"encoding/json"
	"log/slog"

	"bbs-go/internal/pkg/simple/common/strs"
)

func Parse(str string, t interface{}) error {
	if strs.IsAnyBlank(str) {
		return nil
	}
	return json.Unmarshal([]byte(str), t)
}

func ParseBytes(bytes []byte, t interface{}) error {
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

func ToStr(t interface{}) (string, error) {
	if t == nil {
		return "", nil
	}
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ToJsonStr(t interface{}) string {
	if t == nil {
		return ""
	}
	str, err := ToStr(t)
	if err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
		return ""
	}
	return str
}
