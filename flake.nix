{
  description = "A Nix-flake-based Go 1.24 + OCAML development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1.*.tar.gz";

  outputs = { self, nixpkgs }:
    let
      goVersion = 24; # Change this to update the whole stack

      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        };
      });
    in
    {
      overlays.default = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      devShells = forEachSupportedSystem ({ pkgs }: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            # - Go toolchain -
            # go (version is specified by overlay)
            go
            gotools
            golangci-lint
            gopls
            gomodifytags
            gotests
            gore
            godef

            # - OCaml toolchain -
            ocaml  # OCaml compiler
            opam # OPAM package manager
            dune_2 # Dune build system

            # - Supporting tools - 
            protobuf # protoc compiler

            # - C-level libraries for OCaml gRPC & Protobuf plugins - 
            pkg-config
            grpc
            protobufc
            libffi
            zlib
            openssl
          ];

          shellHook = ''
            # 1) Bootstrap OPAM non-interactively
            if ! opam root >/dev/null 2>&1; then
              opam init --bare --disable-sandboxing --no-setup -yes
            fi

            # 2) Create an "empty" local switch that picks up Nix' `ocaml`
            if [ ! -d .opam-switch ]; then
              opam switch create . --empty --yes
            fi

            # 3) Load OPAM environment
            eval #(opam env)

            # 4) Install OCaml dependencies if they're missing
            if ! ocaml list --installed | grep -q grpc-lwc; then
              opam install --yes yojson lwt grpc-lwt grpc ocaml-protoc-plugin
            fi

            # 5) Make sure ocaml-protoc-plugin (and any other OPAM tools) are on PATH
            export PATH=$PATH:$(opam var bin)
          '';
        };
      });
    };
}
