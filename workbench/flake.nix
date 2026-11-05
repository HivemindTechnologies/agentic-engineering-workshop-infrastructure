{
  description = "Agentic engineering workshop, infrastructure track: local k3d cluster, GitOps via Argo CD, Claude Code in a nono sandbox";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = {nixpkgs, ...}: let
    # nono is built for Linux and aarch64-darwin only, so Intel Macs are out.
    systems = ["x86_64-linux" "aarch64-linux" "aarch64-darwin"];
    forEachSystem = f:
      nixpkgs.lib.genAttrs systems (system:
        f (import nixpkgs {
          inherit system;
          config.allowUnfreePredicate = pkg: nixpkgs.lib.getName pkg == "claude-code";
        }));
  in {
    devShells = forEachSystem (pkgs: let
      # Argo CD's own install manifest, fetched by hash: every participant
      # bootstraps byte-identical Argo CD, and a moved upstream file fails the
      # shell instead of the workshop.
      argocdInstall = pkgs.fetchurl {
        url = "https://raw.githubusercontent.com/argoproj/argo-cd/v3.5.5/manifests/install.yaml";
        hash = "sha256-kNyQlL9MWWSv2cDjXXNeHC/LvaT442eZwJLZhk7WX1A=";
      };
      # Shared by both shells; each adds its own kubeconfig.
      commonHook = ''
        export ARGOCD_INSTALL=${argocdInstall}
        # Tool caches inside the repo: the sandbox may write here and nowhere
        # else in $HOME. GOTOOLCHAIN=local forbids silent toolchain downloads.
        # The justfile puts its temp files in .cache/ and fails if it is
        # missing, which it is on a fresh clone.
        mkdir -p "$PWD/.cache"
        export GOPATH=$PWD/.cache/go
        export GOCACHE=$PWD/.cache/go-build
        export GOTOOLCHAIN=local
        # kubectl's discovery cache, by default in ~/.kube/cache.
        export KUBECACHEDIR=$PWD/.cache/kube
        # mkShellNoCC has no C compiler, and the images build without cgo too.
        export CGO_ENABLED=0
      '';
    in rec {
      default = pkgs.mkShellNoCC {
        # Host dependency, on purpose: Docker, because k3d runs its nodes as
        # containers of the host's daemon.
        packages = [
          # macOS ships bash 3.2; just's shebang recipes need bash 5.
          pkgs.bash
          pkgs.git
          # git push runs ssh with -F /dev/null, so no host ssh config is read.
          pkgs.openssh
          pkgs.just
          pkgs.k3d
          pkgs.kubectl
          pkgs.kustomize
          pkgs.kubeconform
          pkgs.kube-linter
          pkgs.go
          pkgs.curl
          pkgs.jq
          pkgs.yq-go
          pkgs.nono
          # Secrets exercise: SOPS files encrypted with a workshop-only age key.
          pkgs.sops
          pkgs.age
          pkgs.claude-code
        ];
        shellHook =
          commonHook
          + ''
            # The read-only agent kubeconfig, written by `just up`. Humans
            # and agents see the cluster through the same view; writes go
            # through git.
            export KUBECONFIG=$PWD/.kubeconfig
          '';
      };
      # `nix develop .#admin`: the default shell with your credentials, for
      # the whole terminal. Same directory as `credentials` in the justfile.
      # Not a boundary against the agent: every host shell can read these
      # files, nono.json is the boundary. It guards against habit: `kubectl`
      # stays read-only unless you opened this shell on purpose, and the
      # audit log (`just audit`) shows its requests as system:admin.
      # Start `just claude` from the default shell: the sandbox cannot read
      # the admin kubeconfig, so its kubectl would fail.
      admin = default.overrideAttrs {
        shellHook =
          commonHook
          + ''
            state=$HOME/.local/state/agentic-workshop-infrastructure
            export KUBECONFIG=$state/admin.kubeconfig
            # No PS1 marker: prompt tools such as starship, sourced from
            # ~/.bashrc, overwrite PS1 before every prompt.
            # Only the exception speaks; on stderr, so `nix develop .#admin
            # -c cmd | jq` still gets clean output.
            echo "Admin shell: kubectl runs as cluster-admin. Type exit to leave." >&2
          '';
      };
    });
  };
}
