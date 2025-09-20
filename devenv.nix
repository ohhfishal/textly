{ pkgs, inputs, lib, config, ... }:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  packages = [
    pkgs-unstable.go
  ];

  git-hooks.hooks = {
    # Shell
    shellcheck.enable = true;

    # Golang
    govet.enable = true;
    gotest.enable = true;
    gofmt.enable = true;
  };
}
