## SimplePAM

a really basic implementation of a Privileged Access Management 

## Setup
`git clone https://github.com/RaynardGerraldo/SimplePAM ; cd SimplePAM ; go mod tidy ; go build`

## Usage

`go run api/main.go api/endpoint.go` -- run this on a seperate terminal

`./SimplePAM admin init` -- initialize admin and server (currently only your localhost)

`./SimplePAM admin add-user <name>` -- add first user for ssh access to server

`./SimplePAM user <user's name>` -- login to allowed server (currently only server-prod)
