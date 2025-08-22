## Foorrent

A simple CLI-based BitTorrent client

## TODO
 - [ ] Handle multiple trackers
 - [ ] Handle DHT and Peer Exchange (PEX)
 - [ ] Better UI
 - [ ] Assign a name for the output file/directory on download
 - [ ] Add support for magnet links
 - [ ] Add support for seeding
 - [ ] Notify/Signal when download is complete

## Known Issues
- Too slow and suspiciously doesn't connect to most peers

## Installation

### Requirements
 - Must have at least Go 1.24

### Steps

 - If you have `make` installed, you can run `make` command to build. Otherwise, run `go build main.go`. You can provide a filename and directory by passing it to `go build` command with `-o` flag.
