INSERT INTO account_groups (name, parent_name, display_after, is_immutable) VALUES
  ('Assets', NULL, NULL, true),
  ('Liabilities', NULL, 'Assets', true),
  ('Equity', NULL, 'Liabilities', true),
  ('Revenues', 'Equity', NULL, true),
  ('Expenses', 'Equity', 'Revenues', true);

INSERT INTO account_types (name, display_after) VALUES
  ('Asset', NULL),
  ('Liability', 'Asset'),
  ('Equity', 'Liability'),
  ('Revenue', 'Equity'),
  ('Expense', 'Revenue');
