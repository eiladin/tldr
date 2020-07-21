<div align="center">
  
  <h1 align="center">TLDR</h1>

  ![build](https://github.com/eiladin/tldr/workflows/build/badge.svg)
  ![test](https://github.com/eiladin/tldr/workflows/test/badge.svg)
  ![GitHub release (latest by date)](https://img.shields.io/github/v/release/eiladin/tldr)
  ![GitHub All Releases](https://img.shields.io/github/downloads/eiladin/tldr/total?color=orange&style=flat)
  ![GitHub last commit](https://img.shields.io/github/last-commit/eiladin/tldr?color=%23ff69b4)
  ![GitHub](https://img.shields.io/github/license/eiladin/tldr?color=blue)
</div>

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