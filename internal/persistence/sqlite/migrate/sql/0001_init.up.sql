CREATE TABLE IF NOT EXISTS account_groups (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	parent_id TEXT REFERENCES account_groups(id),
	display_after TEXT REFERENCES account_groups(id),
	is_immutable BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS account_types (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  display_after TEXT REFERENCES account_types(id)
);

CREATE TABLE IF NOT EXISTS accounts (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  parent_group_id TEXT NOT NULL REFERENCES account_groups(id),
  account_type_id TEXT NOT NULL REFERENCES account_types(id),
  display_after TEXT REFERENCES accounts(id),
  normal_balance TEXT NOT NULL CHECK (normal_balance IN ('debit', 'credit'))
);

CREATE TABLE IF NOT EXISTS journal_entries (
  id TEXT PRIMARY KEY,
  timestamp TEXT NOT NULL,
  description TEXT,
  cross_reference TEXT REFERENCES journal_entries(id) 
);

CREATE TABLE IF NOT EXISTS journal_lines (
  id TEXT PRIMARY KEY,
  account_id TEXT NOT NULL REFERENCES accounts(id),
  amount REAL NOT NULL,
  side TEXT NOT NULL CHECK (side IN ('debit', 'credit')),
  journal_entry_id TEXT NOT NULL REFERENCES journal_entries(id)
);