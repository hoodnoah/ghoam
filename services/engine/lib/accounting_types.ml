type amount = int64 (* represents cents, e.g. $1.23 -> 123L*)

type entry_side = 
  | Debit
  | Credit

type account_type = 
  | Asset
  | ContraAsset
  | Liability
  | Equity
  | Revenue
  | Expense

type normal_balance = entry_side

type account_group = {
  name : string;
  parent_name : string option;
  display_after : account_group option;
  is_immutable : bool;
}

type account = {
  name : string;
  parent_group : account_group option;
  account_type : account_type;
  normal_balance : normal_balance;
  display_after : account option;
}

type journal_entry_line = {
  account : account;
  amount : amount;
  side : entry_side;
}

type journal_entry = {
  id : string;
  timestamp : Ptime.t;
  description : string;
  lines : journal_entry_line list;
  cross_reference : journal_entry option; 
}



