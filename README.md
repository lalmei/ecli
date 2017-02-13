# Ecli

`ecli` is a portable command line client to access the [Keen Eye](https://www.keeneyetechnologies.com/)
digital pathonolgy platform. Its primary use is to ease image upload to
the platform and get slide information in a convenient way.

This [Go](https://golang.org) implementation is currently a work in progress.

## Download

Using `ecli` is easy since it comes as a static binary with no dependencies. Just [grab the latest compiled version](https://github.com/keeneyetech/ecli/releases/latest)
for your system (or compile from source, see below).

## Usage

Just type `ecli` in a terminal to get a list of available commands.

```
Ecli is a command line client for the Keen Eye API. Its primary use is to
perform slide upload and get slide information in a convenient way.

Usage:
  ecli [command]

Available Commands:
  applications List all registered applications
  find         Find slides or groups by criteria
  group        Manage groups
  imageformats List all supported image formats
  label        Manage labels
  login        Open a session
  logout       Close current session
  slide        Manage slides
  version      Show tool version

Flags:
      --config string   Config file (default is $HOME/.ecli.json)
  -q, --quiet           Quiet mode, no verbose output
  -t, --toggle          Help message for toggle

Use "ecli [command] --help" for more information about a command.
```

Tip: for every available commands, you can always get some help and examples by running `ecli [command] -h`.

## Quick Tour

### Opening a Session

TBW

### Manage Labels

Labels provide an easy way to categorize the images and groups based on descriptive
titles. They can have a color and a description.

Any number of labels can be applied to images and groups.

Type `ecli label -h` to see availabe subcommands and examples.

## Compile From Source

First ensure you have a Go 1.5+ working installation then install the traditional way:

    $ go get github.com/keeneyetech/ecli

## License

This is free software, see LICENSE.
