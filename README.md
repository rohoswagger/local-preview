# local-preview

[![Test status](https://github.com/rohoswagger/local-preview/actions/workflows/test.yml/badge.svg)](https://github.com/rohoswagger/local-preview/actions)
[![Deploy Status](https://github.com/rohoswagger/local-preview/actions/workflows/release.yml/badge.svg)](https://github.com/rohoswagger/local-preview/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/jmelahman/local-preview.svg)](https://pkg.go.dev/github.com/jmelahman/local-preview)

A local-first preview-deployment orchestrator: every commit of a registered
git repository becomes a servable preview at its own subdomain
(`<sha>-<repo>.preview.localhost:8080`), built once, deduplicated by
content, and served from a single binary.

`local-preview` resolves each commit into a *(frontend artifact, backend
artifact)* pair by content-addressing the relevant subtrees. Commits that
don't touch the frontend (or backend) reuse the existing artifact — and the
already-running backend process. Backend state follows git lineage: a new
backend version forks its data from the nearest deployed ancestor, so
previews feel continuous on a straight line of commits while divergent
branches can never corrupt each other.

## Install

Download the `preview` binary for your platform from the
[fork releases](https://github.com/rohoswagger/local-preview/releases)
and put it on your `PATH`.

This fork publishes server images to `ghcr.io/rohoswagger/local-preview`.
See [installation](docs/guide/install.md) for image access and release steps.
The Go module path remains `github.com/jmelahman/local-preview`.

## Run

```sh
preview serve
```

The server listens on `:8080`. Open <http://localhost:8080/> for the
dashboard.

Or with [Docker Compose](docs/guide/install.md):

```sh
docker compose up -d
```

## 📖 Documentation

Fork documentation lives in [`docs/`](docs/):

- [Quickstart](docs/guide/quickstart.md) — register a repo, deploy a commit, open the preview.
- [Configuration](docs/guide/configuration.md) — flags and env vars.
- [REST API](docs/reference/api.md) — endpoints exposed by the running server.
- [CLI](docs/reference/cli.md) — `serve`, `repo`, `deploy`, `install-hook`.
