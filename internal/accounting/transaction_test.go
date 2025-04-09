package accounting_test

import (
	"testing"
	"time"

	// module under test
	"github.com/hoodnoah/ghoam/internal/accounting"
)

func TestTransactionBalanced(t *testing.T) {
	t.Run("Balanced Entry", func(t *testing.T) {
		balancedJe := accounting.JournalEntry{
			ID:          "1",
			Timestamp:   time.Now(),
			Description: "Balanced Entry",
			Lines: []accounting.JournalEntryLine{
				{AccountName: "1", Amount: 100.0, Side: accounting.Debit},
				{AccountName: "2", Amount: 100.0, Side: accounting.Credit},
			},
		}

		if !accounting.IsBalanced(balancedJe) {
			t.Errorf("Expected balanced journal entry to be balanced")
		}
	})

	t.Run("Unbalanced Entry", func(t *testing.T) {
		unbalancedJe := accounting.JournalEntry{
			ID:          "2",
			Timestamp:   time.Now(),
			Description: "Unbalanced Entry",
			Lines: []accounting.JournalEntryLine{
				{AccountName: "1", Amount: 100.0, Side: accounting.Debit},
				{AccountName: "2", Amount: 50.0, Side: accounting.Credit},
			},
		}

		if accounting.IsBalanced(unbalancedJe) {
			t.Errorf("Expected unbalanced journal entry to be unbalanced")
		}
	})

	t.Run("Zero Entry", func(t *testing.T) {
		zeroJe := accounting.JournalEntry{
			ID:          "3",
			Timestamp:   time.Now(),
			Description: "Zero Entry",
			Lines:       []accounting.JournalEntryLine{},
		}

		if !accounting.IsBalanced(zeroJe) {
			t.Errorf("Expected zero journal entry to be balanced")
		}
	})
}
