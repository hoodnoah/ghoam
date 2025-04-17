package sqlite

import (
	// std
	"context"
	"database/sql"

	// external

	_ "github.com/mattn/go-sqlite3" // sqlite driver

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

type accountGroupRepo struct {
	db *sql.DB
}

type validateFn func(context.Context, *accounting.AccountGroup) error

// gets an account group by name.
//
// Returns ErrGroupNotFound if the account does not exist.
func (r *accountGroupRepo) GetByName(ctx context.Context, name string) (accounting.AccountGroup, error) {
	const query = `
		SELECT name, parent_name, display_after, is_immutable
	  FROM account_groups
		WHERE name = ?;
	`

	var group accounting.AccountGroup
	err := r.db.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&group.Name,
		&group.ParentName,
		&group.DisplayAfter,
		&group.IsImmutable,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return accounting.AccountGroup{}, &accounting.ErrGroupNotFound{Name: name}
		}
		return accounting.AccountGroup{}, err
	}

	return group, nil
}

// Upsert inserts or updates an AccountGroup in the database
// More or less an upsert
func (r *accountGroupRepo) Upsert(ctx context.Context, group *accounting.AccountGroup) error {
	const query = `
		INSERT INTO account_groups
		  (name, parent_name, display_after, is_immutable)
		VALUES
		  (?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			parent_name = excluded.parent_name,
			display_after = excluded.display_after,
			is_immutable = excluded.is_immutable;
	`

	// check if the group exists already
	extantGroup, err := r.GetByName(ctx, group.Name)
	if err != nil && !accounting.IsGroupNotFound(err) {
		return err
	}

	// if the group is immutable, it cannot be updated.
	if extantGroup.IsImmutable {
		return &accounting.ErrGroupImmutable{Name: group.Name}
	}

	// insert/update if the record either doesn't exist, or is mutable
	_, err = r.db.ExecContext(
		ctx,
		query,
		group.Name,
		group.ParentName,
		group.DisplayAfter,
		group.IsImmutable,
	)

	return err
}

// Gets all account groups
func (r *accountGroupRepo) GetAll(ctx context.Context) ([]*accounting.AccountGroup, error) {
	const query = `
		SELECT name, parent_name, display_after, is_immutable
		FROM account_groups;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var groups []*accounting.AccountGroup

	for rows.Next() {
		var group accounting.AccountGroup

		if err := rows.Scan(
			&group.Name,
			&group.ParentName,
			&group.DisplayAfter,
			&group.IsImmutable,
		); err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	// sort accountGroups in-place
	if err := accounting.SortAccountGroupsInPlace(groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *accountGroupRepo) Insert(ctx context.Context, group *accounting.AccountGroup) error {
	const query = `
		INSERT INTO account_groups
		  (name, parent_name, display_after, is_immutable)
		VALUES
		  (?, ?, ?, ?);
	`

	// run pre-insert validations
	err := r.runValidators(ctx, group, r.validateGroupExists, r.validateParentExists, r.validateDisplayAfterExists)
	if err != nil {
		return err
	}

	// we know the error is ErrAccountNotFound - we can insert
	_, err = r.db.ExecContext(ctx, query,
		group.Name,
		group.ParentName,
		group.DisplayAfter,
		false,
	)

	return err
}

// utility method for determining if a name exists
func (r *accountGroupRepo) nameExists(ctx context.Context, name string) (bool, error) {
	_, err := r.GetByName(ctx, name)

	// no error means it's been found
	if err == nil {
		return true, nil
	}

	// if it's anything other than a groupNotFound, the error should be propagated
	if !accounting.IsGroupNotFound(err) {
		return false, err
	}

	return false, nil
}

// determines if a group's name exists
func (r *accountGroupRepo) validateGroupExists(ctx context.Context, group *accounting.AccountGroup) error {
	groupNameExists, err := r.nameExists(ctx, group.Name)
	if err != nil {
		return err
	}

	if groupNameExists {
		return &accounting.ErrGroupAlreadyExists{Name: group.Name}
	}

	return nil
}

// determines if a parent exists
func (r *accountGroupRepo) validateParentExists(ctx context.Context, group *accounting.AccountGroup) error {
	parentGroupExists, err := r.nameExists(ctx, group.ParentName.String)
	if err != nil {
		return err
	}
	if !parentGroupExists {
		return &accounting.ErrParentNameNotExists{Name: group.ParentName.String}
	}

	return nil
}

// determines if a referenced displayAfter exists
func (r *accountGroupRepo) validateDisplayAfterExists(ctx context.Context, group *accounting.AccountGroup) error {
	if group.DisplayAfter.Valid {
		afterGroupExists, err := r.nameExists(ctx, group.DisplayAfter.String)
		if err != nil {
			return err
		}

		if !afterGroupExists {
			return &accounting.ErrDisplayAfterNameNotExists{Name: group.DisplayAfter.String}
		}
	}

	return nil
}

// variadic validation running utility method
func (r *accountGroupRepo) runValidators(ctx context.Context, group *accounting.AccountGroup, validationFns ...validateFn) error {
	for _, validateFn := range validationFns {
		err := validateFn(ctx, group)
		if err != nil {
			return err
		}
	}

	return nil
}
