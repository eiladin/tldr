# TLDR

Go implementation of a cli for https://github.com/tldr-pages/tldr

run `tldr {command}` in a terminal where `{command}` is a standard cli command (e.g. `tldr curl`)

* Caches to .tldr by downloading a zip from https://github.com/tldr-pages/tldr
* Outputs color coded output to stdout
* Detects OS for loading appropriate syntax and examples
* Help and argument parsing implemented via kingpin (https://github.com/alecthomas/kingpin)
* Allows override of the platform via the `--platform` argument
* Allows for fetching a random page via the `--random` argument

## Releasing

Releases are done through goreleaser (https://goreleaser.com)

If UPX is installed, then the binaries will be compressed

### Step-by-step
* Add tag (git tag -a v1.0.0 -m "message")
* push tag (git push origin v1.0.0)
* make publish

