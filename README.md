## Foorrent
A simple CLI-based BitTorrent client

## Motivation
This project is created as a way to deepen my understanding of the [BitTorrent Protocol](https://en.wikipedia.org/wiki/BitTorrent). I wanted to explore how the peer-to-peer and distributed aspect of such systems works in practice, and building this project allowed me to do that. The goal of this project is to implement a simple client that runs in the terminal where you can pass it a `.torrent` file and it handles contacting the tracker, setting up the connection with the peers, requesting pieces and blocks, downloading the pieces, and saving it to disk.

## Why the name?
I always struggle with naming my projects and fallback to naming it foo, bar, foobar. Thus, the name comes from "foo" and "torrent" combined (who would've guessed lol🤣)

## Scope
To limit the scope of this project and realistically create an output in a reasonable time some major aspects of the client and the protocol is not handled. Major features not included in the scope includes the following:

 - Handling multiple trackers
 - Integration of DHT and Peer Exchange (PEX)
 - Support for magnet links
 - Seeding support

## What's Next?
 - [ ] Handle DHT and Peer Exchange (PEX)
 - [ ] Better UI
 - [ ] Assign a name for the output file/directory on download
 - [ ] Add support for magnet links
 - [ ] Add support for seeding
 - [ ] Notify/Signal when download is complete

## Known issues worth investigating
 - Doesn't/fails to connect to most of the peers returned by the tracker.
 - Download speed is slow.

## Installation

### Pre-requisites
As a part of this learning experience, I deliberately avoided relying on external packages whenever possible. Because of this approach, the only requirement for building and running this in your machine is to at least have the version of `Go` I used in this project, `Go 1.24.2`. Visit [Go website](https://go.dev) for download instructions.

### Building
If you have `make` installed or your running this in a `mac` or `linux` environment, simply run:
```bash
make
```

Building without using make:
```bash
go build main.go
```
This will build the executable as `main`. You can specify the name with the `-o` flag in go build.

## Usage
View the `help` instructions to see the available commands.
```bash
foorrent -h
# or
foorrent --help
```
