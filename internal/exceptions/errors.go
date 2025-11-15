package exceptions

import "fmt"

type MiddlewareFatalError struct {
	desc string
}

func (err *MiddlewareFatalError) Error() string {
	return fmt.Sprintf("middleware fatal error %v :", err.desc)
}