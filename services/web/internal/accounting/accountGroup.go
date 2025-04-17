package accounting

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// utility type for a validator function
type groupValidatorFn func(*AccountGroup) error

// these are the base account groups
var reservedAccountGroupNames = []string{
	"Assets",
	"Liabilities",
	"Equity",
	"Revenue",
	"Expenses",
}

// representation of an account group
type AccountGroup struct {
	Name         string         `json:"name"`
	ParentName   sql.NullString `json:"parent_name"`   // applicable for a subheading; empty -> null
	DisplayAfter sql.NullString `json:"display_after"` // for orderings; empty -> null
	IsImmutable  bool           `json:"is_immutable"`  // if true, this group cannot be deleted or altered
}

// utility method for deep equality
func (a *AccountGroup) Equals(other *AccountGroup) bool {
	return a.Name == other.Name && a.ParentName == other.ParentName && a.DisplayAfter == other.DisplayAfter && a.IsImmutable == other.IsImmutable
}

// constructor for a new AccountGroup
//
// parentName is provided as a non-nullable string; all account groups other than
// the base account groups must have a parent.
//
// DisplayAfter is still nullable, since some AccountGroups might not have a predecessor
// within their parent.
//
// These account groups cannot be made immutable, as well. Only the default account groups are immutable.
func NewAccountGroup(
	name string,
	parentName string,
	displayAfter sql.NullString,
) (*AccountGroup, error) {
	// set up an unvalidated account group
	newGroup := AccountGroup{}
	newGroup.Name = name
	newGroup.ParentName.Valid = true
	newGroup.ParentName.String = parentName
	newGroup.DisplayAfter = displayAfter
	newGroup.IsImmutable = false

	// validate the accountgroup
	err := runValidators(&newGroup,
		validateGroupName,
		validateGroupNameNotReserved,
		validateGroupParentName,
		validateDisplayAfter,
	)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

// determines if a given AccountGroup.Name is a valid string
func validateGroupName(group *AccountGroup) error {
	if strings.EqualFold(group.Name, "") {
		return errors.New("AccountGroup requires a non-empty string for its Name attribute")
	}
	return nil
}

// determines if a given AccountGroup.Name is a reserved name
func validateGroupNameNotReserved(group *AccountGroup) error {
	for _, r := range reservedAccountGroupNames {
		if strings.EqualFold(group.Name, r) {
			return fmt.Errorf("received a reserved AccountGroupName: %s", group.Name)
		}
	}

	return nil
}

// determines if a given AccountGroup.ParentName is a valid string
func validateGroupParentName(group *AccountGroup) error {
	if !group.ParentName.Valid {
		return errors.New("AccountGroup requires its ParentName.Valid == true")
	}

	if group.ParentName.String == "" {
		return errors.New("AccountGroup requires its ParentName.String to be non-empty")
	}

	return nil
}

// determines if a given AccountGroup.DisplayAfter is valid
func validateDisplayAfter(group *AccountGroup) error {
	// intentionally empty
	if !group.DisplayAfter.Valid && group.DisplayAfter.String != "" {
		return fmt.Errorf("expected a DisplayAfter with Valid == false to contain an empty string, received %s", group.DisplayAfter.String)
	}

	if group.DisplayAfter.Valid && group.DisplayAfter.String == "" {
		return errors.New("expected a DisplayAfter with Valid == true to have a non-empty String value")
	}

	return nil
}

// runs a variadic number of validators against an AccountGroup.
func runValidators(group *AccountGroup, validators ...groupValidatorFn) error {
	for _, validatorFn := range validators {
		if err := validatorFn(group); err != nil {
			return err
		}
	}

	return nil
}
