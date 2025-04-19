module QC = QCheck
module Gen = QCheck.Gen
module QR = QCheck_runner
module PC = Ptime_clock
module Types = Engine.Accounting_types
module Journal_entry = Engine.Journal_entry

(* dummy account – fields other than “amount”/“side” don’t matter here *)
let dummy_account : Types.account = {
  name           = "X";
  parent_group   = None;
  account_type   = Asset;
  normal_balance = Debit;
  display_after  = None;
}

(* generator of a single line with a random positive amount & side *)
let gen_amt : int64 Gen.t =
  Gen.map Int64.of_int (Gen.int_range 1 10_000)  (* 1¢ to 100 $ of cents *)

let gen_side : Types.entry_side Gen.t =
  Gen.oneofl [Types.Debit; Types.Credit]

let gen_line : Types.journal_entry_line Gen.t =
  Gen.(
    gen_amt  >>= fun amt ->
    gen_side >>= fun side ->
    return { Types.account=dummy_account; amount=amt; side }
  )

(* generator for a journal_entry with exactly one line *)
let gen_single_line_je : Types.journal_entry Gen.t = 
  Gen.map(fun line -> {
    Types.id = "single";
    timestamp = (PC.now ());
    description = "single-line test";
    lines = [line];
    cross_reference = None
  }) gen_line

let arb_single_line_je = QCheck.make gen_single_line_je

(* property: any single‑line entry is unbalanced *)
let prop_single_line_never_balanced =
  QC.Test.make
    ~count:1_000
    ~name:"single‑line entry never balanced"
    arb_single_line_je
    (fun je -> not (Journal_entry.is_balanced je))

(* run it *)
let () =
  QCheck_runner.run_tests_main
    [ prop_single_line_never_balanced ]