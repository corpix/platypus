with import <nixpkgs>{};
{ pkgs ? import <nixpkgs> {} }:

buildGo19Package rec {
  name = "platypus-unstable-${version}";
  version = "development";

  buildInputs = with pkgs; [ git glide ];

  src = ./.;
  goPackagePath = "github.com/cryptounicorns/platypus";
}
