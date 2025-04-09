package accounting

import "fmt"

type ErrGroupNotFound struct {
	Name string
}

type ErrGroupImmutable struct {
	Name string
}

func (e *ErrGroupNotFound) Error() string {
	return fmt.Sprintf("account group \"%s\" not found", e.Name)
}

func (e *ErrGroupImmutable) Error() string {
	return fmt.Sprintf("account group \"%s\" exists already and is immutable.", e.Name)
}

// helper utilities
func IsGroupNotFound(err error) bool {
	_, ok := err.(*ErrGroupNotFound)
	return ok
}

func IsGroupImmutable(err error) bool {
	_, ok := err.(*ErrGroupImmutable)
	return ok
}
