# shell.nix
{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  # Tools you want available in the dev shell
  buildInputs = with pkgs; [
    go          # Go toolchain
    gopls       # Language server (for editors/IDE)
    delve       # Debugger
    go-tools    # misc go tools (vet, etc)
  ];

  # Optional: set up GOPATH and basic env
  shellHook = ''
    export GOPATH=$PWD/.gopath
    export GOBIN=$GOPATH/bin
    export PATH=$GOBIN:$PATH
    echo "GOPATH set to $GOPATH"
  '';
}

