{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-25.05";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { ... } @ inputs: inputs.utils.lib.eachDefaultSystem (system:
    let
      pkgs = import inputs.nixpkgs { inherit system; };
    in
    {
      devShell = pkgs.mkShell {
        packages = [ pkgs.go ];
        shellHook = ''
          export GOPATH="$HOME/go"
          export PATH="$PATH:$GOPATH/bin"
        '';
      };
    });
}
