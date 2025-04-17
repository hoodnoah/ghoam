CREATE TABLE IF NOT EXISTS account_groups (
	name TEXT PRIMARY KEY,
	parent_name TEXT REFERENCES account_groups(name),
	display_after TEXT REFERENCES account_groups(name),
	is_immutable BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS account_types (
  name TEXT PRIMARY KEY,
  display_after TEXT REFERENCES account_types(name)
);

CREATE TABLE IF NOT EXISTS accounts (
  name TEXT PRIMARY KEY,
  parent_group_name TEXT NOT NULL REFERENCES account_groups(name),
  account_type TEXT NOT NULL REFERENCES account_types(name),
  display_after TEXT REFERENCES accounts(name),
  normal_balance TEXT NOT NULL CHECK (normal_balance IN ('Debit', 'Credit'))
);

CREATE TABLE IF NOT EXISTS journal_entries (
  id TEXT PRIMARY KEY,
  timestamp TEXT NOT NULL,
  description TEXT,
  cross_reference TEXT REFERENCES journal_entries(id) 
);

CREATE TABLE IF NOT EXISTS journal_lines (
  id TEXT PRIMARY KEY,
  account_name TEXT NOT NULL REFERENCES accounts(name),
  amount REAL NOT NULL,
  side TEXT NOT NULL CHECK (side IN ('Debit', 'Credit')),
  journal_entry_id TEXT NOT NULL REFERENCES journal_entries(id)
);