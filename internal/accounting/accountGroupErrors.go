package accounting

import "fmt"

type ErrGroupNotFound struct {
	Name string
}

type ErrGroupImmutable struct {
	Name string
}

type ErrGroupAlreadyExists struct {
	Name string
}

type ErrParentNameNotExists struct {
	Name string
}

type ErrDisplayAfterNameNotExists struct {
	Name string
}

func (e *ErrGroupNotFound) Error() string {
	return fmt.Sprintf("account group \"%s\" not found", e.Name)
}

func (e *ErrGroupImmutable) Error() string {
	return fmt.Sprintf("account group \"%s\" exists already and is immutable.", e.Name)
}

func (e *ErrGroupAlreadyExists) Error() string {
	return fmt.Sprintf("account group \"%s\" already exists", e.Name)
}

func (e *ErrParentNameNotExists) Error() string {
	return fmt.Sprintf("account group parent \"%s\" does not exist", e.Name)
}

func (e *ErrDisplayAfterNameNotExists) Error() string {
	return fmt.Sprintf("account group display after \"%s\" does not exist", e.Name)
}

// --------- helper utilities ------------
func IsGroupNotFound(err error) bool {
	_, ok := err.(*ErrGroupNotFound)
	return ok
}

func IsGroupImmutable(err error) bool {
	_, ok := err.(*ErrGroupImmutable)
	return ok
}

func IsGroupAlreadyExists(err error) bool {
	_, ok := err.(*ErrGroupAlreadyExists)
	return ok
}

func IsParentNameNotExists(err error) bool {
	_, ok := err.(*ErrParentNameNotExists)
	return ok
}

func IsDisplayAfterNotExists(err error) bool {
	_, ok := err.(*ErrDisplayAfterNameNotExists)
	return ok
}
