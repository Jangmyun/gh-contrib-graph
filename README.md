# gh-contrib-graph

A GitHub CLI extension to display your contribution graph directly in the terminal.

![screenshot](docs/screenshot.gif)

Design inspired by [ghstats-cli](https://github.com/d1rshan/ghstats-cli).

## Install

```sh
gh extension install Jangmyun/gh-contrib-graph
```

Requires `gh auth login` to already be set up — the extension reuses your existing `gh` session, no separate token needed.

## Usage

```sh
# your own contribution graph
gh contrib-graph

# someone else's
gh contrib-graph octocat
```
