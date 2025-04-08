INSERT INTO account_groups (id, name, parent_id, display_after, is_immutable) VALUES
  ('assets', 'Assets', NULL, NULL, true),
  ('liabilities', 'Liabilities', NULL, 'assets', true),
  ('equity', 'Equity', NULL, 'liabilities', true),
  ('revenues', 'Revenues', 'equity', NULL, true),
  ('expenses', 'Expenses', 'equity', 'revenues', true);

INSERT INTO account_types (id, name, display_after) VALUES
  ('asset', 'Asset', NULL),
  ('liability', 'Liability', 'asset'),
  ('equity', 'Equity', 'liability'),
  ('revenue', 'Revenue', 'equity'),
  ('expense', 'Expense', 'revenue');
