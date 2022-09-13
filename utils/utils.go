package utils

import (
	"encoding/json"
	"os"
	"path"
)

func ReadJSON(filePath string, v interface{}) error {
	cwd, _ := os.Getwd()
	buffer, err := os.ReadFile(path.Join(cwd, filePath))
	if err != nil {
		return err
	}

	err = json.Unmarshal(buffer, v)
	return err
}

func WriteJSON(filePath string, v interface{}) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	buffer, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(cwd, filePath), buffer, 0666)
	return err
}
