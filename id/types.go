package id

import (
	"fmt"
	"strconv"
)

type ID int64

func (curr ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", curr)), nil
}

func (curr *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	bytes := data
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	v, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return err
	}
	*curr = ID(v)
	return nil
}

func (curr ID) String() string {
	return strconv.FormatInt(int64(curr), 10)
}

func Parse(str string) (ID, error) {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(v), nil
}
