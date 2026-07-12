package syntax

import (
	"errors"
	"strconv"
)

func errorBadLength(name string, length int) error {
	return errors.New(`Invalid length '` + strconv.FormatInt(int64(length), 10) + `' for ` + name)
}

func errorBadType(name string) error {
	return errors.New(`Incompatible input type for ` + name)
}
