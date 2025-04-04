package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins bro", r)
	return []byte(strconv.Quote(jsonValue)), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	var strJson string = string(jsonValue)
	// Just in case for debug the string: %#v\n
	strJson, err := strconv.Unquote(strJson)

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(strJson, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}
	// conv to int
	mins, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(mins)

	return nil
}
