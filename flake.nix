{
  description = "A Nix-flake-based Go 1.24 + OCAML development environment";

  inputs = {
    # List of platform identifiers, e.g. "x86_64-linux" etc.
    systems.url = "github:nix-systems/default"; 

    # Snapshot of nixpkgs, pinned by a FlakeHub wildcard.
    nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";
  };


  # ──────────────────────────────────────────────────────────
  # outputs : receives materialized inputs and *returns* an attr‑set
  # ──────────────────────────────────────────────────────────
  outputs = {self, nixpkgs, systems}:
    let
      lib = nixpkgs.lib; # Nixpkgs standard library
      eachSystem = lib.genAttrs (import systems);
    in
    {
      packages = eachSystem (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          # Build ./services/web as a Go module
          web = pkgs.buildGoModule {
            pname = "web";
            version = "0.1.0";

            # Source derivation relative to *this flake*
            src = ./services/web;

            # Tell Nix to fetch and vendor all module dependencies.
            # The first time, use a dummy SHA-256. Then copy and paste the real hash from the build error.
            vendorHash = "sha256-ojT4jqhBXaPHZD80aiZIkL2A/cBnbPVaoOO+J3g22WY=";

            # We want Go 1.24; the pkg set already contains go_1_24
            buildInputs = [pkgs."go_1_24"];
          };

          # `nix build` with no name falls back to building web
          default = self.packages.${system}.web;
        }
      );

      devShells = eachSystem (system:
        let 
          pkgs = nixpkgs.legacyPackages.${system};
          go = pkgs.go_1_24;
        in
        {
          default = pkgs.mkShell {
            # packages placed on $PATH
            packages = with pkgs; [
              # --- Go toolchain ---
              go
              gotools
              golangci-lint
              gopls
              gomodifytags
              gotests
              godef
            ];

            # Expose everything that the 'web' derivation builds with
            inputsFrom = [self.packages.${system}.web];
          };
        }
      );

      checks = eachSystem (system:
        let 
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          # re-use the build definition, but leave only the test phase enabled
          web-tests = self.packages.${system}.web.overrideAttrs (old: {
            name = "test-${old.pname}";

            installPhase = ''
              mkdir -p $out
            '';
          });
        }
      );
    };
}
