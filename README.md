# oauth2-proxy-k8s-impersonation

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![GitHub Release](https://img.shields.io/github/v/release/aslafy-z/oauth2-proxy-k8s-impersonation)](https://github.com/aslafy-z/oauth2-proxy-k8s-impersonation)
[![Go Reference](https://pkg.go.dev/badge/github.com/golang-templates/seed.svg)](https://pkg.go.dev/github.com/aslafy-z/oauth2-proxy-k8s-impersonation)
[![go.mod](https://img.shields.io/github/go-mod/go-version/aslafy-z/oauth2-proxy-k8s-impersonation)](go.mod)
[![LICENSE](https://img.shields.io/github/license/aslafy-z/oauth2-proxy-k8s-impersonation)](LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/aslafy-z/oauth2-proxy-k8s-impersonation/build)](https://github.com/aslafy-z/oauth2-proxy-k8s-impersonation/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/aslafy-z/oauth2-proxy-k8s-impersonation)](https://goreportcard.com/report/github.com/aslafy-z/oauth2-proxy-k8s-impersonation)
[![Codecov](https://codecov.io/gh/aslafy-z/oauth2-proxy-k8s-impersonation/branch/main/graph/badge.svg)](https://codecov.io/gh/aslafy-z/oauth2-proxy-k8s-impersonation)

This is a tool that proxies requests to Kubernetes Dashboard by adding an Authorization header with a token from a service account and maps oauth2-proxy response headers to Kubernetes impersonation headers.

## Usage

TBD

## Setup

TBD

## Build

### Terminal

- `make` - execute the build pipeline.
- `make help` - print help for the [Make targets](Makefile).

### Visual Studio Code

`F1` → `Tasks: Run Build Task (Ctrl+Shift+B or ⇧⌘B)` to execute the build pipeline.

## Release

The release workflow is triggered each time a tag with `v` prefix is pushed.

_CAUTION_: Make sure to understand the consequences before you bump the major version. More info: [Go Wiki](https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher), [Go Blog](https://blog.golang.org/v2-go-modules).

## Maintainance

Remember to update Go version in [.github/workflows](.github/workflows), [Makefile](Makefile) and [devcontainer.json](.devcontainer/devcontainer.json).

Notable files:

- [devcontainer.json](.devcontainer/devcontainer.json) - Visual Studio Code Remote Container configuration,
- [.github/workflows](.github/workflows) - GitHub Actions workflows,
- [.github/dependabot.yml](.github/dependabot.yml) - Dependabot configuration,
- [.vscode](.vscode) - Visual Studio Code configuration files,
- [.golangci.yml](.golangci.yml) - golangci-lint configuration,
- [.goreleaser.yml](.goreleaser.yml) - GoReleaser configuration,
- [Dockerfile](Dockerfile) - Dockerfile used by GoReleaser to create a container image,
- [Makefile](Makefile) - Make targets used for development, [CI build](.github/workflows) and [.vscode/tasks.json](.vscode/tasks.json),
- [go.mod](go.mod) - [Go module definition](https://github.com/golang/go/wiki/Modules#gomod),
- [tools.go](tools.go) - [build tools](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module).

## Contributing

Simply create an issue or a pull request.

## Credits

This repository uses the [golang-templates/seed](https://github.com/golang-templates/seed) repository as a template.
