# Go Packet Playground

## Description

This is a repo to implement a network parser based on Golang.

## Pcap source

- [WireShark pcap samples](https://wiki.wireshark.org/SampleCaptures)

```bash
make
make clean
make help
```

## Dependancies

```
go get github.com/google/gopacket@v1.1.17
```

## Plan

- Complet core package and HTTP plugin
- Document on core
- MySQL plugin
- Output interface with Protocol buffer message
- MSSQL, PostgreSQL, MongoDB, Redis
