{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/23.05";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      myShell = { buildInputs ? [], shellHook ? ""}:
      let
        gotainer = self.outputs.packages.${system}.gotainer;
        underlyingShell = pkgs.mkShell {
          buildInputs = buildInputs ++ [ pkgs.bashInteractive  ];
          inherit shellHook;
        };
      in
        pkgs.mkShell {
          buildInputs = [ gotainer ];
          shellHook = ''
          ${gotainer}/bin/gotainer run  ${underlyingShell} ${pkgs.bashInteractive}
          exit
          '';
        }
     ;
      in
    {

      packages.${system} = {
        hello = myShell {
          buildInputs = [
            pkgs.snakemake
          ];
        };
        gotainer = pkgs.buildGoModule {
            name = "gotainer";
            version = "0.0";
            src = ./gotainer;
            vendorSha256 = null;
       };
      };

      devShells.${system} = {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # 
            bashInteractive
            go
          ];
          shellHook = ''
            echo ${pkgs.bashInteractive}
          '';
        };
      };

    };
}
