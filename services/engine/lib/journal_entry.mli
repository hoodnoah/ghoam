open Accounting_types

val sum_debit_lines : journal_entry -> int64
val sum_credit_lines : journal_entry -> int64
val is_balanced : journal_entry -> bool