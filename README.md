# Ecli

`ecli` is a portable command line client to access the [Keen Eye](https://www.keeneyetechnologies.com/)
digital pathonolgy platform. Its primary use is to ease image upload to
the platform and get slide information in a convenient way.

This [Go](https://golang.org) implementation is currently a work in progress.

## See it In Action

Please have a quick look at this 2 minutes video showing some useful commands to upload a slide to the Keen Eye platform:

[![asciicast](https://asciinema.org/a/103308.png)](https://asciinema.org/a/103308)

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

### Creating a Config File

`ecli` needs a config file to run. It holds profile information about the credentials to
be used in order to connect and use the Keen Eye API.

By default, `ecli` will look for an `.ecli.json` file in the current user's home directory but you may want to specify a different file by using the `--config` flag.

The config file contains profile information in a JSON file like
```
{
  "profile1": {
    "login": "<YOUR LOGIN EMAIL ADDRESS>",
    "password": "<YOUR PASSWORD>",
    "url": "https://<YOUR URL TO PLAFORM>/api/v2"
  }
}
```

with `<YOUR URL TO PLAFORM>` being the address used by your browser to connect to the platform, like `https://prefix.keeneyetechnologies.com`. Note that the `url` defined in the config file must end with `/api/v2`.

### Sessions

Before you can start using the API with `ecli`, you need to open a session. Log in to the service with `ecli login`:
```
$ ecli login profile1
profile1: you have been logged in successfully.
```
Now you can start using the API with `ecli`.

Your session may expire 30 minutes later due to inactivity. In this case, you will have to login again. You can also
explicitely log out from the service:
```
$ ecli logout
Your session is now closed.
```

### Labels

Labels provide an easy way to categorize the images and groups based on descriptive
titles. They can have a color and a description.

Any number of labels can be applied to images and groups.

Type `ecli label -h` to see availabe subcommands and examples.

## Groups

Groups are like directories in a filesystem, useful to organize images. They can have
a name and a description. Labels can be added to groups.

Type `ecli group -h` to see availabe subcommands and examples.

## Image Formats

Use the `imageformats` command to list all supported image formats for your installation. The first column is the image format value to pass to the `slide upload` command during an image upload (see section below).
```
$ ecli imageformats
tiff                           tiff                           TIFF image
ndpi                           ndpi                           NDPI image
ndpis                          ndpis                          NDPIS image
dicom                          dicom                          DICOM image
```
This is the result from a default installation.

## Slide Upload

`ecli slide upload` can take many parameters, like image format, pixel size (value and unit) and so on. Please type `ecli slide upload -h` to get the full list of available options.

For example, uploading a TIFF image with a 1 micron pixel size and apply two "retina" and "core" labels on it can be performed with
```
$ ecli slide upload stained_cells.tif -f tiff -p 1 -l "retina" -l "core"
```

By default, the slide will be uploaded to the root of the work list. When uploading a slide inside a group somewhere else in the hierarchy is required, the `--group-id` flag can be used.

Say you need to upload a slide into the existing `Study 1` group. You first need to get the group unique ID to use during the upload:
```
$ ecli group ls
58a173aee779892825712398 "Other Group"                
58a43cc5e779890c2486ca6d "Study 1"
```
Group `Study 1` has ID `58a43cc5e779890c2486ca6d`, so to upload the slide into that group, just
```
$ ecli slide upload stained_cells.tif --group-id 58a43cc5e779890c2486ca6d
```

## Compile From Source

First ensure you have a Go 1.5+ working installation then install the traditional way:

    $ go get github.com/keeneyetech/ecli

## License

This is free software, see LICENSE.
