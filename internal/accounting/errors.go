package accounting

import "fmt"

type ErrAccountNotFound struct {
	Name string
}

func (e *ErrAccountNotFound) Error() string {
	return fmt.Sprintf("account \"%s\" not found", e.Name)
}

// helper utility
func IsAccountNotFound(err error) bool {
	_, ok := err.(*ErrAccountNotFound)
	return ok
}
