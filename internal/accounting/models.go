package accounting

import "time"

// enumeration of the types of accounts
type AccountType string

const (
	Asset       AccountType = "asset"
	ContraAsset AccountType = "contra_asset"
	Liability   AccountType = "liability"
	Equity      AccountType = "equity"
)

// enumeration of the types of balances
type NormalBalance string

const (
	CreditNormal NormalBalance = "credit"
	DebitNormal  NormalBalance = "debit"
)

type Account struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	ParentID      *string       `json:"parent_id"`
	Type          AccountType   `json:"type"`
	NormalBalance NormalBalance `json:"normal_balance"`
}

// enumeration of the types of entries
type EntrySide string

const (
	Debit  EntrySide = "debit"
	Credit EntrySide = "credit"
)

// representation of a single journal entry line
type JournalEntryLine struct {
	AccountID string    `json:"account_id"`
	Amount    float64   `json:"amount"`
	Side      EntrySide `json:"side"`
}

// representation of a journal entry, comprised of
// multiple journal entry lines
type JournalEntry struct {
	ID          string             `json:"id"`
	Timestamp   time.Time          `json:"timestamp"`
	Description string             `json:"description"`
	Lines       []JournalEntryLine `json:"lines"`
}
