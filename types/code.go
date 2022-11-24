package types

import "fmt"

type Code struct {
	Code    uint
	Message string
}

func (curr *Code) Error() string {
	return fmt.Sprintf("code:%v,message:%v", curr.Code, curr.Message)
}
