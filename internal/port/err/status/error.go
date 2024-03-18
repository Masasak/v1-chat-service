package status

import "strconv"

type Error struct {
	Msg  string
	Code int
}

func (e Error) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Msg
}

func NewErr(status int, msg string) *Error {
	return &Error{Code: status, Msg: msg}
}
