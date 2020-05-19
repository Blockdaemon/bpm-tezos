# tezos

This is a package for the Blockchain Package Manager by [Blockdaemon](https://blockdaemon.com/). Deploy, maintain, and upgrade blockchain nodes on your own infrastructure.

Further reading:

* End-user documentation: https://cli.bpm.docs.blockdaemon.com/
* Developer documentation: https://sdk.bpm.docs.blockdaemon.com/

# Contributing

Pleaes use [conventional commits](https://www.conventionalcommits.org) for you commit messages. This will help us in the future to auto-generate changelogs.

New features should be developed in a branch and merged after a code review.

# Building from source

## Requirements

Make sure you have the following tools:

- [Go](https://golang.org/) is the main pogramming language. It needs to be installed
- [goreleaser](https://goreleaser.com/) is used to build binary packages. It needs to be installed
- [golangci-lint](https://github.com/golangci/golangci-lint) is used to do static code checks. It needs to be installed
- [GPG](https://gnupg.org/) is used to sign build artifacts. It needs to be installed

## Building during development

For quick development builds for your own system run:

    make build

To install the development package run:

    bpm packages install --from-file ./bin/bpm-tezos-development

## Building a test release:

To build a test release for multiple platforms without publishing it run:

    make test-release

## Releasing a new version

You need the GPG key to sign the releases imported into your GPG keyring.

    make VERSION=<version> GITLAB_TOKEN=<gitlab-token> release

`<version>` needs to be a valid [semantic version](https://semver.org/). Do **not prefix with `v`**, the script does that automatically.

`<gitlab-token>` needs to be a gitlab token with api scope. You can create a token here: https://gitlab.com/profile/personal_access_tokens.

The artifacts will be published in Gitlab. Contact support@blockdaemon.com to add your package to the Blockchain Package Registry.