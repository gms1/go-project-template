# go-project-template

[![GitHub License](https://img.shields.io/github/license/gms1/go-project-template)](https://raw.githubusercontent.com/gms1/go-project-template/refs/heads/main/LICENSE)
[![build Status](https://github.com/gms1/go-project-template/actions/workflows/build.yaml/badge.svg)](https://github.com/gms1/go-project-template/actions/workflows/build.yaml)
[![Coverage Status](https://codecov.io/gh/gms1/go-project-template/branch/main/graph/badge.svg)](https://codecov.io/gh/gms1/go-project-template)
[![Go Report Card](https://goreportcard.com/badge/github.com/gms1/go-project-template)](https://goreportcard.com/report/github.com/gms1/go-project-template)
[![Go Reference](https://pkg.go.dev/badge/github.com/gms1/go-project-template?status.svg)](https://pkg.go.dev/github.com/gms1/go-project-template?tab=doc)
[![Release](https://img.shields.io/github/release/gms1/go-project-template.svg?style=flat-square)](https://github.com/gms1/go-project-template/releases)

## usage

- install prerequisites

  - mandatory: gofumpt, golangci-lint and goreleaser

  ```bash
  go install mvdan.cc/gofumpt@latest
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  go install github.com/goreleaser/goreleaser/v2@latest
  ```

  - optional and recommended: pre-commit

  ```bash
  pip install pre-commit
  ```

- export this template to a new project folder

  ```bash
  ./scripts/export.sh <target-project-name> [target-project-owner]
  ```

- in new project folder

  - to format, lint, build and test:

  ```bash
  ./scripts/all.sh
  ```

  - verify, commit and push everything to the new git repository

  - optional install pre-commit hook

  ```bash
  pre-commit install
  pre-commit install --hook-type commit-msg
  ```

  - create a release

  ```bash
  ./scripts/release.sh <version>
  ```

  this updates the `Version` in "pkg/common/about.go", generates the docs in "docs" and commits and pushes the changes using the commit message "release: <version>".

  A "v<version>" tag is then created and pushed, which triggers the release workflow from .github/workflows/release.yaml.

  This release workflow then creates the binaries and archives for the configured OS/Arch combinations, as well as a multi-platform docker image. If the later is not needed, the “Dockerfile” can be removed, as well as the "docker" and "docker_manifests" sections in ".goreleaser.yaml"
