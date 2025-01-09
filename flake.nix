{
  description = "Dynamic IP Address Watcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    allSystems = [
      "x86_64-linux"
      "aarch64-linux"
      "x86_64-darwin"
      "aarch64-darwin"
    ];

    forAllSystems = f:
      nixpkgs.lib.genAttrs allSystems (system:
        f {
          pkgs = import nixpkgs {
            inherit system;
            overlays = [self.overlays.default];
            config.allowUnfree = true;
          };
          inherit system;
        });
  in {
    overlays.default = final: prev: let
    in {
      dynamic-ip-watcher = final.buildGoModule rec {
        pname = "dynamic-ip-watcher";
        version = "0.0.1";
        src = final.lib.cleanSource self;
        vendorHash = "sha256-/xfMDF5sL7U3f+WLT+xE31WtIHdL3VZmL977DFY7PPM=";
        ldflags = ["-s" "-w" "-X main.version=v${version}"];
        outputName = "dynamic-ip-watcher";
      };
    };

    packages = forAllSystems ({
      pkgs,
      system,
      ...
    }: {
      default = pkgs.dynamic-ip-watcher;
    });

    nixosModules.dynamic-ip-watcher = import ./nix/module.nix self.overlays.default;

    devShells = forAllSystems ({
      pkgs,
      system,
      ...
    }: {
      default = pkgs.mkShell {
        packages = with pkgs; [
          go_1_23
          gotools
          go-mockery
          golangci-lint
          direnv
        ];
      };
    });
  };
}
