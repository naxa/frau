# popuko

[![Build Status (master)](https://travis-ci.org/student-kyushu/frau.svg?branch=master)](https://travis-ci.org/student-kyushu/frau)
[![CircleCI](https://circleci.com/gh/student-kyushu/frau/tree/master.svg?style=svg)](https://circleci.com/gh/student-kyushu/frau/tree/master)

## What is this?

**This is fork repository. Original is [popuko](https://github.com/voyagegroup/popuko).**

frau is a very cute and useful GitHub bot.

She can do as follow:

* merge a pull request automatically
* assign someone to review a pull request
* check if a pull request is mergeable

## Motivation

To just go for it is student rights, we decided to introduce bots and service to our project.

frau is one of our projects and trying. She makes our development process useful. At first, we were going to use bors, or popuko as it is, however, we wanted more experience and incentive. And we thought some features we needed weren't suitable for popuko's policy. Therefore, we forked popuko and started to develop.

## Features

Almost features are based on popuko, see [here](https://github.com/voyagegroup/popuko#features). Only the features we added are shown.

* parse one by one line, not only first line
* replace `me` with `reviewer` when used `r=me`
* read not just only comments but also description of a pull request
* see label `S-do-not-merge`

## Setup Instructions

### Build & Launch the Application

1. This requires that [`github.com/golang/dep`](https://github.com/golang/dep) has been installed.
2. Build from the source.
    * Run these steps:
      1. `make bootstrap`.
      2. `make build` or `make build_linux_x64`.
    * Run `make help` to see more details.
    * You also can do `go get`.
3. Create the config directory.
    * By default, this app uses `$XDG_CONFIG_HOME/frau/` as the config dir.
      (If you don't set `$XDG_CONFIG_HOME` environment variable, this use `~/.config/frau/`.)
    * You can configure the config directory by `--config-base-dir`
4. Set `config.toml` to the config directory.
    * Let's copy from [`./example.config.toml`](./example.config.toml)
5. Start the exec binary.
    * This app dumps all logs into stdout & stderr.
    * If you'd like to use TLS, then provide `--tls`, `--cert`, and `--key` options.

#### Set up for your repository in GitHub

1. Set the account (or the team which it belonging to) which this app uses as a collaborator
   for your repository (requires __write__ priviledge).
2. Add `OWNERS.json` file to the root of your repository.
    * Please see [`OwnersFile`](./setting/ownersfile.go) about the detail.
    * The example is [here](./OWNERS.json).
3. Set `http://<your_server_with_port>/github` for the webhook to your repository with these events:
    * `Issue comment`
    * `Push`
    * `Status` (required to use Auto-Merging feature (non GitHub App CI services)).
    * `Check Suite` (required to use Auto-Merging feature (GitHub App CI Services)).
    * `Pull Request` (required to remove all status (`S-` prefixed) labels after a pull request is closed).
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

[The MIT License](https://github.com/student-kyushu/frau/blob/master/LICENSE.MIT)

## How to Contribute

We welcome your contributing, thanks!

See [CONTRIBUTING.md](https://github.com/student-kyushu/frau/blob/master/CONTRIBUTING.md)

[homu]: https://github.com/barosl/homu
[servo-homu]: https://github.com/servo/homu
[highfive]: https://github.com/servo/highfive
[bors.tech]: https://bors.tech/
[github-rust-repo]: https://github.com/rust-lang/
[github-servo]: https://github.com/servo
[graydon's-entry]: http://graydon2.dreamwidth.org/1597.html
[bors-ng]: https://github.com/bors-ng/bors-ng
