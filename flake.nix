{
  description = "nix dev env config";

  # Flake inputs
  inputs = {
    nixpkgs.url = "nixpkgs/22.11";
    nixpkgs-unstable.url = "nixpkgs/nixpkgs-unstable";
  };

  # Flake outputs
  outputs = { self, nixpkgs, nixpkgs-unstable }:
    let
      # Systems supported
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit Apple Silicon
      ];

      # Helper to provide system-specific attributes
      nameValuePair = name: value: { inherit name value; };
      genAttrs = names: f: builtins.listToAttrs (map (n: nameValuePair n (f n)) names);
      forAllSystems = f: genAttrs allSystems (system: f {
        pkgs = import nixpkgs {
          inherit system;
          #config = { allowUnfree = true; };
        };
        pkgsUnstable = import nixpkgs-unstable {
          inherit system;
          #config = { allowUnfree = true; };
        };
      });
    in
    {
      # Development environment output
      devShells = forAllSystems ({ pkgs, pkgsUnstable }: {
        # The default development shell
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            kind
            ctlptl
            tilt
            k9s
          ] ++ (with pkgsUnstable; [
            # something in unstable
          ]);
        };
      });
    };
}
