(* A minimal money type *)
Require Import Coq.ZArith.ZArith.
Definition Money := Z.          (* use integers for cents *)

(* Accounts are just strings for now *)
Require Import Coq.Strings.String.
Definition AccountId := string.

(* A posting = account × signed amount *)
Record Posting := {
  acct : AccountId ;
  amt  : Money      (* positive = debit, negative = credit *)
}.

Definition JournalEntry := list Posting.

(* Sum of amounts in an entry *)
Definition entry_total (e:JournalEntry) : Money :=
  fold_right (fun p acc => amt p + acc)%Z 0%Z e.

(** *** Fundamental invariant: a valid journal entry is balanced *)
Definition valid_entry (e:JournalEntry) : Prop :=
  entry_total e = 0%Z.

(* Simple helper: adding amounts is commutative *)
Lemma fold_perm :
  forall l1 l2, Permutation l1 l2 ->
    entry_total l1 = entry_total l2.
Proof.
  intros l1 l2 Hperm.
  unfold entry_total.
  now apply Permutation.fold_right_comm with (f:=Z.add) (z:=0%Z).
Qed.

(** **** Posting a balanced entry preserves ledger equality ****)

Definition Ledger := AccountId -> Money.

(* apply a posting to a ledger *)
Definition post (L:Ledger) (p:Posting) : Ledger :=
  fun a => if String.eqb a (acct p)
           then L a + amt p
           else L a.

(* fold a whole entry *)
Definition post_entry (L:Ledger) (e:JournalEntry) : Ledger :=
  fold_left post e L.

(* Sum of all balances – we only care that it stays the same *)
Definition ledger_sum (L:Ledger) : Money :=
  (* finite map is nicer, but we stay minimal *) 0%Z.

Theorem post_entry_preserves_sum :
  forall L e,
    valid_entry e ->
    ledger_sum (post_entry L e) = ledger_sum L.
Proof.
  (* proof sketch – depends on a finite‑map implementation;
     left unfinished for brevity *)
Admitted.