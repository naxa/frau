# frau

[![Build Status](https://github.com/naxa-jp/frau/workflows/Go%20CI/badge.svg)](https://github.com/naxa-jp/frau/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/naxa-jp/frau)](https://goreportcard.com/report/github.com/naxa-jp/frau)

## What is this?

**This is fork repository. Original is [popuko](https://github.com/voyagegroup/popuko).**

frau is a very cute and useful GitHub bot.

She can do as follows:

* merge a pull request automatically
* assign someone to review a pull request
* check if a pull request is mergeable

## Motivation

To just go for it is student rights, we decided to introduce bots and services to our projects.

frau is one of our projects and tries. She makes our development process useful. At first, we were going to use bors, or popuko as it is, however, we wanted more experience and incentive. And we thought some features we needed weren't suitable for popuko's policy. Therefore, we forked popuko and started to develop.

## Features

Almost features are based on popuko, see [popuko's features](https://github.com/voyagegroup/popuko#features) and [frau's CHANGELOG](https://github.com/naxa-jp/frau/blob/master/CHANGELOG.md).

## Setup Instructions

### Build & Launch the Application

Please get binary from [here](https://github.com/naxa-jp/frau/releases) or build by yourself.

#### Build process

1. This requires that [`Go`](https://github.com/golang/go) and [`Git`](https://git-scm.com/) have been installed.
2. Build from the source. Run these steps:
    1. `git clone https://github.com/naxa-jp/frau.git`
    2. `cd frau && go build`

#### Launch process

1. Create the config directory.
    * By default, this app uses `$XDG_CONFIG_HOME/frau/` as the config dir.
        (If you don't set `$XDG_CONFIG_HOME` environment variable, this use `~/.config/frau/`.)
    * You can configure the config directory by `--config-base-dir`
2. Set `config.toml` to the config directory.
    * Let's copy from [`./example.config.toml`](./example.config.toml)
3. Start the exec binary.
    * This app dumps all logs into stdout & stderr.
    * If you'd like to use TLS, then provide `--tls`, `--cert`, and `--key` options.

#### Set up for your repository in GitHub

1. Set the account (or the team which it belonging to) which this app uses as a collaborator
    for your repository (requires __write__ priviledge).
2. Add `OWNERS.json` file to the root of your repository.
    * Please see [`OwnersFile`](./setting/ownersfile.go) about the detail.
    * The example is [here](./OWNERS.json).
3. Set `http://<your_server_with_port>/github` for the webhook to your repository with these events:
    * `Check suites` (required to use Auto-Merging feature (GitHub App CI Services))
    * `Issue comments`
    * `Pull requests` (required to remove all status (`S-` prefixed) labels after a pull request is closed)
    * `Pull request reviews` (required to `r+` PRs via reviews)
    * `Pushes`
    * `Statuses` (required to use Auto-Merging feature (non GitHub App CI services))
4. Create these labels to make the status visible.
    * `S-awaiting-review`
        * for a pull request assigned to some reviewer.
    * `S-awaiting-merge`
        * for a pull request queued to this bot.
    * `S-needs-rebase`
        * for an unmergeable pull request.
    * `S-fails-tests-with-upstream`
        * for a pull request which fails tests after try to merge into upstream (used by Auto-Merging feature).
    * `S-do-not-merge`
        * for a pull request not to want to merge
5. Enable to start the build on creating the branch named `auto` for your CI service (e.g. TravisCI).
    * You can configure this branch's name by `OWNERS.json`.
6. Done!

## License

[The MIT License](https://github.com/naxa-jp/frau/blob/master/LICENSE.MIT)

## How to Contribute

We welcome your contributing, thanks!

See [CONTRIBUTING.md](https://github.com/naxa-jp/frau/blob/master/CONTRIBUTING.md)
