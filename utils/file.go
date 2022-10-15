package utils

import (
	"os"
	"path"
)

func ReadFile(filePath string) string {
	cwd, err := os.Getwd()
	CatchError(Fatal, err)

	buffer, err := os.ReadFile(path.Join(cwd, filePath))
	CatchError(Fatal, err)

	return string(buffer)
}

func WriteFile(filePath string, buffer string) {
	cwd, err := os.Getwd()
	CatchError(Fatal, err)

	outPath := path.Join(cwd, filePath)
	os.MkdirAll(path.Dir(outPath), 0666)

	err = os.WriteFile(outPath, []byte(buffer), 0666)
	CatchError(Fatal, err)
}
