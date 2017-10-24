# droplist

Go application that configures firewall to block IP addresses from [Spamhaus DROP list](https://www.spamhaus.org/drop/drop.txt)

## Outstanding Items
1) Clear previously run entries before running next batch (this will make utility idempotent)
2) Provide SystemD service definition

## Build Instructions

1) Linux AMD64: `env GOOS=linux GOARCH=amd64 go build`

## Download

Binary download for Linux AMD64: [Download](https://github.com/mkez00/droplist/raw/master/resources/droplist.zip)

## Vagrant

Execute `vagrant up` in directory to show demo of utility (Ubuntu 16.04).