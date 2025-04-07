package accounting

import (
	"database/sql"
	"time"
)

// representation of an account group
type AccountGroup struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	ParentID     sql.NullString `json:"parent_id"`     // applicable for a subheading; empty -> null
	DisplayAfter sql.NullString `json:"display_after"` // for orderings; empty -> null
	IsImmutable  bool           `json:"is_immutable"`  // if true, this group cannot be deleted or altered
}

// enumeration of the types of accounts
type AccountType string

const (
	Asset       AccountType = "asset"
	ContraAsset AccountType = "contra_asset"
	Liability   AccountType = "liability"
	Equity      AccountType = "equity"
	Revenue     AccountType = "revenue"
	Expense     AccountType = "expense"
)

// enumeration of the types of balances
type NormalBalance string

const (
	CreditNormal NormalBalance = "credit"
	DebitNormal  NormalBalance = "debit"
)

type Account struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	ParentGroupID string         `json:"account_group_id"`
	AccountType   AccountType    `json:"account_type"`
	NormalBalance NormalBalance  `json:"normal_balance"`
	DisplayAfter  sql.NullString `json:"display_after"` // for ordering; empty -> null
}

// enumeration of the types of entries
type EntrySide string

const (
	Debit  EntrySide = "debit"
	Credit EntrySide = "credit"
)

// representation of a single journal entry line
type JournalEntryLine struct {
	AccountID      string         `json:"account_id"`
	Amount         float64        `json:"amount"`
	Side           EntrySide      `json:"side"`
	CrossReference sql.NullString `json:"cross_reference"` // e.g. in a reversal; empty -> null
}

// representation of a journal entry, comprised of
// multiple journal entry lines
type JournalEntry struct {
	ID          string             `json:"id"`
	Timestamp   time.Time          `json:"timestamp"`
	Description string             `json:"description"`
	Lines       []JournalEntryLine `json:"lines"`
}
