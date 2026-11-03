# Agentic Engineering Workshop: Infrastructure, Preparation

Welcome, and thanks for joining! We look forward to the workshop with you.
Please complete the setup below beforehand, so we can start right away.

It downloads every tool and image the workshop needs, so the workshop day
needs no downloads: about 2 GB compressed, about 8 GB on disk once unpacked.

## Setup

We use the package manager **Nix**; `flake.nix` declares all tools.
Nix stores them in `/nix/store`, separate from your other software.
They are only on your `PATH` inside the `nix develop` shell.

The workshop runs a small Kubernetes cluster on your laptop, as Docker
containers. Linux and macOS on Apple silicon are supported.

1. Install Nix. We recommend the Determinate installer (Linux and macOS, flakes enabled):
   https://docs.determinate.systems/determinate-nix/
   If you use the official installer (https://nixos.org/download/), enable flakes:
   `mkdir -p ~/.config/nix && echo "experimental-features = nix-command flakes" >> ~/.config/nix/nix.conf`
2. Install and start Docker: Docker Engine on Linux
   (https://docs.docker.com/engine/install/), Docker Desktop
   (https://docs.docker.com/desktop/) or OrbStack (https://orbstack.dev/) on macOS.
3. In this repo, enter the dev shell: `nix develop`
   This opens a shell with the tools from `flake.nix`. `exit` leaves it.
4. Check it works: `claude --version && k3d version && docker info`
5. Download the images: `just prepare`
   It starts a cluster, fills its image cache, a Docker volume named
   `infra-workshop-image-cache`, and deletes the cluster again; the cache
   stays for the workshop. If it stops halfway, run `just prepare` again: it
   continues where it stopped.

## After the Workshop

Remove the image cache: `docker volume rm infra-workshop-image-cache`

To uninstall Nix, follow https://manual.determinate.systems/installation/uninstall.html (Determinate installer) or https://nix.dev/manual/nix/stable/installation/uninstall.html (official installer).
