{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/23.05";
    nsc.url = "github:GuilloteauQ/nix-shell-container";
  };

  outputs = { self, nixpkgs, nsc }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system} = {
        my-shell = nsc.lib.mkShell {
          buildInputs = with pkgs; [
            python3
          ];
        };
      };
    };
}
