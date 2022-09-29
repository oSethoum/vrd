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

	outPath := path.Join(cwd, filePath)
	os.MkdirAll(path.Dir(outPath), 0666)

	err = os.WriteFile(outPath, buffer, 0666)
	return err
}
