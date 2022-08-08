# gotorrent

[![version](https://img.shields.io/github/v/release/ismaelpadilla/gotorrent)](https://github.com/ismaelpadilla/gotorrent/releases)
[![golangci-lint](https://github.com/ismaelpadilla/gotorrent/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ismaelpadilla/gotorrent/actions/workflows/golangci-lint.yml)


TUI for searching torrents. You can open a torrent's magnet link in your default app, or download its .torrent file. This app does not handle leeching/seeding a torrent.

Currently queries ThePirateBay's API.

https://user-images.githubusercontent.com/7772501/180335527-d8a9678f-8e61-429d-bbc3-1a085884059d.mp4

# Installation

## Using go
```sh
go install github.com/ismaelpadilla/gotorrent@latest
```

Alternatively, you can clone the repo:

```sh
git clone https://github.com/ismaelpadilla/gotorrent
cd gotorrent
go install
```

### Install from package

#### AUR

gotorrent is available on the AUR:
```
yay -S gotorrent
```

#### Homebrew

```
brew install ismaelpadilla/tap/gotorrent
```

# Usage:

```sh
gotorrent <query>
```

Input a number and press enter to navigate to that torrent's magnet link. Or use the `up` and `down` (or `j`/`k`) keys to navigate the torrent list.

## Keybinds

- `up`/`k`: Scroll up.
- `down`/`j`: Scroll down.
- `Enter`: Navigate to a selected torrent.
- `t`: Download .torrent file.
- `c`: Copy magnet link to clipboard.
- `d`: See torrent description.
- `f`: See torrent files.
- `s`: Enter a new search query.
- `q`: Quit.
- `?`: Expand/minimize help.

## Flags

```
  -d, --debug                    show debug information
  -f, --download-folder string   folder where files are downloaded
  -h, --help                     help for gotorrent
  -p, --persist                  keep gotorrent open after selecting torrent
```

# Configuration

A configuration file can be used, but flags take precedence over configuration.

## Location

A `config` file can be put in the following locations:

- Same folder as the executable.
- `$HOME/.config/gotorrent/`.

## Config keys

Only one configuration key can be set:

`download-folder`: Same as the `--download-folder` flag.

## Configuration file example

```toml
download-folder = "/home/myUser/torrent"
```


# Roadmap

Work on the [v0.1.0 milestone](https://github.com/ismaelpadilla/gotorrent/milestone/1) has been completed and the first version has been officially released.

Currenlty work is being done on:
- Code cleanup
- Bugfixes
- QoL things
- Potentially adding new clients to search torrents from other sources
