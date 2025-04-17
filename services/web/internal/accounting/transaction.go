package accounting

// Checks if a journal entry is balanced.
//
// A journal entry is considered balanced if the total
// line debits equal the total line credits.
func IsBalanced(je JournalEntry) bool {
	var debitTotal, creditTotal float64

	for _, line := range je.Lines {
		switch line.Side {
		case Debit:
			debitTotal += line.Amount
		case Credit:
			creditTotal += line.Amount
		}
	}

	return debitTotal == creditTotal
}
