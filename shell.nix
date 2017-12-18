with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "mole-shell";
  buildInputs = [
    ncurses
    go
    gocode
    go-bindata
    glide
    godef
    bison
  ];
}
