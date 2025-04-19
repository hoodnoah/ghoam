let () = 
  let open Alcotest in
  run "Event_sourcing"
    [
      ("basics",
        [
          test_case "tautology" `Quick (fun () ->
            check int "numbers agree" 2 2);
        ]
      )
    ]