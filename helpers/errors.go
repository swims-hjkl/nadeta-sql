package helpers

import (
	"errors"
	"fmt"
)

func GetMissingArgError(arg string) error {
	return errors.New(fmt.Sprintf("--%s is mandatory", arg))
}
