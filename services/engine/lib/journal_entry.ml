open Accounting_types

(** Sum the entries on a given side, debit or credit *)
let sum_lines_by_side 
  (target_side : entry_side)
  (entry : journal_entry)
  : int64 = 
  List.fold_left
    (fun accumulator {amount; side; _} ->
      if side = target_side
      then Int64.add accumulator amount
    else accumulator)
  0L
  entry.lines

(* partial application for summing debit lines *)
let sum_debit_lines : journal_entry -> int64 = 
  sum_lines_by_side Debit

(* partial application for summing credit lines *)
let sum_credit_lines : journal_entry -> int64 = 
  sum_lines_by_side Credit

(** Determine if a journal entry is balanced 
  e.g. all debits equal all credits*)
let is_balanced (je : journal_entry) : bool =
  sum_debit_lines je = sum_credit_lines je