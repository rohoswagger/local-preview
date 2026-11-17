# Install

The server manages git repos in-process (no `git` binary needed) and runs
your projects' own build commands, so whatever toolchains your targets build
with — e.g. `go`, `bun` — must be on the `PATH`.

## GitHub releases

Download a prebuilt archive for your platform from the
[releases page](https://github.com/rohoswagger/local-preview/releases) and
place the `preview` binary on your `PATH`.

## Docker

```sh
docker run -d --name preview \
  --restart unless-stopped \
  -p 127.0.0.1:8080:8080 \
  -v preview-data:/data \
  ghcr.io/rohoswagger/local-preview:latest
```

::: warning
The published image bundles no build toolchains (`git`, `go`, `bun`, …),
so on its own it can only build targets whose build commands need nothing
beyond a shell. Builds for repos that declare an `image` in their manifest
(see [Build images](../reference/preview-toml.md#build-images)) can instead
run on the host's Docker daemon — see below. Backend **run** commands still
execute inside the server's container either way, so backends must be
self-contained binaries (or need nothing beyond a shell) when the server
runs in Docker. Run the binary on the host when your targets need more.
:::

### Container builds from Compose

The repo's `compose.yaml` gives the composed server what it needs to hand
`image`-declared build steps to the host's Docker daemon:

- the daemon socket mounted at `/var/run/docker.sock` (source overridable
  via `DOCKER_SOCK_PATH`), and
- the data directory bind-mounted at the **same absolute path** on the
  host and in the container. The daemon resolves bind sources against its
  own filesystem, so the scratch paths the server passes it must mean the
  same thing on both sides — a named volume can't provide that.

The service runs as container root so the default invocation just works —
holding the docker socket is root-equivalent on the daemon anyway:

```sh
docker compose up -d
```

To keep the data dir's files owned by your host user instead, run
non-root and grant the supplementary GID that owns the socket:

```sh
mkdir -p ~/.local/share/local-preview
DOCKER_SOCK_GID=$(stat -c '%g' /var/run/docker.sock) \
PREVIEW_UID=$(id -u) PREVIEW_GID=$(id -g) \
  docker compose up -d
```

Rootless daemon — point at the rootless socket; container root maps to
your host user, so the default user is already right:

```sh
DOCKER_SOCK_PATH=$XDG_RUNTIME_DIR/docker.sock docker compose up -d
```

Repos registered by a local path (rather than a clone URL) must also be
bind-mounted into the container at that same path — see the commented
example in `compose.yaml`.

## From source

```sh
git clone https://github.com/rohoswagger/local-preview
cd local-preview
bun install --cwd web --frozen-lockfile && bun run --cwd web build
go build -tags embed -o preview .
```

## Fork releases

The `rohoswagger/local-preview` fork publishes GitHub release archives and
GHCR images from version tags matching `v*.*.*`. It does not publish to PyPI.
The workflow uses `GITHUB_TOKEN`; no Docker Hub or PyPI credentials are needed.
Keep the Go module path `github.com/jmelahman/local-preview` unchanged.

After the first image publication, open the package settings and set
`local-preview` visibility to **Public**. GHCR initially creates private
packages, even for public repositories. See
[GitHub's container registry documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry).
This one-time step permits anonymous pulls from EC2 hosts. Verify access
with an empty Docker credential directory before deployment:

```sh
docker_config=$(mktemp -d)
DOCKER_CONFIG="$docker_config" docker pull ghcr.io/rohoswagger/local-preview:latest
rm -rf "$docker_config"
```

Use the published image digest in infrastructure. Keep the control service,
ingestion workers, and CI CLI on the same protocol release.
