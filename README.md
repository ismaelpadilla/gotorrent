# gotorrent

TUI for searching torrents. Currently queries ThePirateBay's API.

This is an early work in progress.

# Installation

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

# Todo

- [ ] Improve UI (display an actual table with seeders, leechers, etc).
- [ ] Error handling.
- [ ] Fluff (colors, etc).
