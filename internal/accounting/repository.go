package accounting

import "context"

type AccountGroupRepository interface {
	Save(ctx context.Context, group *AccountGroup) error
	GetByName(ctx context.Context, name string) (AccountGroup, error)
	GetAll(ctx context.Context) ([]*AccountGroup, error)
}

type AccountRepository interface {
	Save(ctx context.Context, account *Account) error
	GetAll(ctx context.Context) ([]*Account, error)
	ByName(ctx context.Context, Name string) (Account, error)
	// ListByGroup(ctx context.Context, groupID string) ([]Account, error)
}

type JournalEntryRepository interface {
	Save(ctx context.Context, je JournalEntry) error
	// ByID(ctx context.Context, id string) (JournalEntry, error)
	// ListByAccount(ctx context.Context, accountID string) ([]JournalEntry, error)
}
