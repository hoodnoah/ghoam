DELETE FROM account_types
  WHERE name IN ("Asset", "Liability", "Equity", "Revenue", "Expense");

DELETE FROM account_groups
  WHERE name IN ("Assets", "Liabilities", "Equity", "Revenues", "Expenses");