{
  description = "A very basic flake";

  inputs = {
    # This one is overridable
    nixpkgs.url = "github:nixos/nixpkgs/23.05";
    # This one, probably not
    nixpkgs_stable.url = "github:nixos/nixpkgs/23.05";
  };

  outputs = { self, nixpkgs, nixpkgs_stable }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      pkgs_stable = import nixpkgs_stable { inherit system; };
    in {
      lib = {
        mkShell = { containerize ? false, ... }@shellInputs:
          let
            gotainer = self.outputs.packages.${system}.gotainer;
            underlyingShell = pkgs.mkShell shellInputs;
          in if containerize then
            pkgs.mkShell {
              buildInputs = [ gotainer ];
              shellHook = ''
                ${gotainer}/bin/gotainer run ${underlyingShell} ${pkgs.bashInteractive}
                exit
              '';
            }
          else
            underlyingShell;
      };

      packages.${system} = {
        hello = self.lib.mkShell {
          buildInputs = [ pkgs.snakemake pkgs.python3 pkgs.lua ];
          containerize = true;
          shellHook = ''
            lua -v
            python3 --version
          '';
        };
        gotainer = pkgs_stable.buildGoModule {
          name = "gotainer";
          version = "0.0";
          src = ./gotainer;
          vendorSha256 = null;
        };
      };

      devShells.${system} = {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [ bashInteractive go ];
          shellHook = ''
            echo ${pkgs.bashInteractive}
          '';
        };
        record = pkgs.mkShell {
          buildInputs = with pkgs; [ vhs ];
        };
      };

    };
}
