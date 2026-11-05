{
  description = "Agentic engineering workshop, infrastructure track: the slides";

  # A flake of its own, so the participants' workbench/ shell carries no slide
  # tooling, and the slide tool can change without touching the workbench.
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = {nixpkgs, ...}: let
    systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin"];
    forEachSystem = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
  in {
    devShells = forEachSystem (pkgs: {
      default = pkgs.mkShellNoCC {
        # The software track's deck tooling (.bin/), with its npm
        # devDependencies taken from nixpkgs at the same versions instead.
        packages = [
          pkgs.just
          pkgs.python3
          pkgs.reveal-md
          pkgs.entr
        ];
        # The version is pinned here; its update check would only nag, and it
        # writes to ~/.config.
        shellHook = ''
          export NO_UPDATE_NOTIFIER=1
        '';
      };
    });
  };
}
