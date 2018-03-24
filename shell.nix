with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "nix-cage-shell";
  buildInputs = [
    nsq
    go
    gocode
    glide
    godef
  ];
  shellHook = ''
    export GOPATH=~/projects
  '';
}
