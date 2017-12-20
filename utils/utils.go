package utils

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func CheckGlError(prefix string) error {
	if errCode := gl.GetError(); errCode != 0 {
		errMessage := fmt.Sprintln(prefix, ":", errCode)
		return errors.New(errMessage)
	}

	return nil
}
