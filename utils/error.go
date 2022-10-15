package utils

import (
	"fmt"
	"os"
)

type ErrorMode = string

const (
	Fatal    ErrorMode = "Error"
	Warninig ErrorMode = "Warining"
)

func CatchError(mode ErrorMode, err error) {
	if err != nil && mode == Fatal {
		fmt.Println("🔴", mode, err.Error())
		os.Exit(0)
	} else if err != nil {
		fmt.Println("🟡", mode, err.Error())
	}
}
