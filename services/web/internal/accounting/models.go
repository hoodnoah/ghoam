package accounting

import (
	"database/sql"
	"time"
)

// enumeration of the types of accounts
type AccountType string

const (
	Asset       AccountType = "Asset"
	ContraAsset AccountType = "Contra Asset"
	Liability   AccountType = "Liability"
	Equity      AccountType = "Equity"
	Revenue     AccountType = "Revenue"
	Expense     AccountType = "Expense"
)

// enumeration of the types of balances
type NormalBalance string

const (
	CreditNormal NormalBalance = "Credit"
	DebitNormal  NormalBalance = "Debit"
)

type Account struct {
	Name            string         `json:"name"`
	ParentGroupName string         `json:"parent_group_name"`
	AccountType     AccountType    `json:"account_type"`
	NormalBalance   NormalBalance  `json:"normal_balance"`
	DisplayAfter    sql.NullString `json:"display_after"` // for ordering; empty -> null
}

// helper functions for sorting
func AccountIDFn(account *Account) string {
	return account.Name
}

func AccountAfterFn(account *Account) (string, bool) {
	if account.DisplayAfter.Valid {
		return account.DisplayAfter.String, true
	}
	return "", false
}

// enumeration of the types of entries
type EntrySide string

const (
	Debit  EntrySide = "Debit"
	Credit EntrySide = "Credit"
)

// representation of a single journal entry line
type JournalEntryLine struct {
	AccountName    string         `json:"account_name"`
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
