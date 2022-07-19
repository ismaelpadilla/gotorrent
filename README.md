# gotorrent

[![golangci-lint](https://github.com/ismaelpadilla/gotorrent/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/ismaelpadilla/gotorrent/actions/workflows/golangci-lint.yml)

TUI for searching torrents. Currently queries ThePirateBay's API.

# Installation

```sh
go install github.com/ismaelpadilla/gotorrent@latest
```

Alternatively, you can clone the repo:

```sh
git clone https://github.com/ismaelpadilla/gotorrent
cd gotorrent
go install
```

# Usage:

```sh
gotorrent <query>
```

Input a number and press enter to navigate to that torrent's magnet link. Or use the `up` and `down` (or `j`/`k`) keys to navigate.

# Roadmap

While basic functionality is already implemented, there is still some work to do before publishing an official release. This is documented in the [v0.1.0 project](https://github.com/users/ismaelpadilla/projects/1/views/1) and the [v0.1.0 milestone](https://github.com/ismaelpadilla/gotorrent/milestone/1).
