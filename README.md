<p align="center">
  
  <h3 align="center">TLDR</h3>

  <p align="center">
    <a href="https://github.com/eiladin/tldr/actions?query=workflow%3Abuild"><img alt="Build" src="https://github.com/eiladin/tldr/workflows/build/badge.svg"></a>
    <a href="https://github.com/eiladin/tldr/actions?query=workflow%3Atest"><img alt="Test" src="https://github.com/eiladin/tldr/workflows/test/badge.svg"></a>
    <a href="https://github.com/eiladin/tldr/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/eiladin/tldr"></a>
    <a href="https://github.com/eiladin/tldr/releases/latest"><img alt="Downloads" src="https://img.shields.io/github/downloads/eiladin/tldr/total?color=orange"></a>
    <a href="https://github.com/eiladin/tldr/tree/checkout-tfe-support"><img alt="Latest Commit" src="https://img.shields.io/github/last-commit/eiladin/tldr?color=ff69b4"></a>
  </p>
</p>

---

## Synopsis

Go implementation of a cli for https://github.com/tldr-pages/tldr


## Usage
run `tldr {command}` in a terminal where `{command}` is a standard cli command (e.g. `tldr curl`)

## Features
* Caches to .tldr by downloading a zip from https://github.com/tldr-pages/tldr
* Outputs color coded output to stdout
* Detects OS for loading appropriate syntax and examples
* Help and argument parsing implemented via kingpin (https://github.com/alecthomas/kingpin)
* Allows override of the platform via the `--platform` argument
* Allows for fetching a random page via the `--random` argument

## Installation
### Homebrew (MacOS and LinuxBrew)

Tap the cask:

`brew tap eiladin/homebrew-tldr`

Install tldr:

`brew install eiladin/homebrew-tldr/tldr` 

**Note:** there is a name-conflict on `tldr` so the full cask name must be used for install

Uninstall tldr:

`brew uninstall tldr`