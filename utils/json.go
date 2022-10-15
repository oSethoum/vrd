package utils

import "encoding/json"

func ReadJson(filePath string, v interface{}) {
	err := json.Unmarshal([]byte(ReadFile(filePath)), v)
	CatchError(Fatal, err)
}

func WriteJson(filePath string, v interface{}) {
	buffer, err := json.Marshal(v)
	CatchError(Fatal, err)
	WriteFile(filePath, string(buffer))
}
