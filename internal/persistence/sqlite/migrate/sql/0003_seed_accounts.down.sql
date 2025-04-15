DELETE FROM accounts
  WHERE name = "Retained Earnings"
  AND parent_group_name = "Equity"
  AND account_type = "Equity"
  AND display_after = NULL
  AND normal_balance = "Credit";