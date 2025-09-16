{
  description = "A dotfile manager for pair programmers and lazy people";
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      version = builtins.substring 0 10 self.rev or "dirty";
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          dfm = pkgs.buildGoModule {
            pname = "dfm";
            inherit version;
            src = ./.;
            vendorHash = "sha256-2xjwWrP1O2brpzNuUxG/CZtMoBiG1V5BOlsp8MSI8WQ=";
            ldflags = [
              "-X main.Version=${version}"
            ];
          };
        });

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.dfm);
    };
}

