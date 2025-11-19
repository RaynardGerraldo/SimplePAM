## SimplePAM

a really basic implementation of a Privileged Access Management 

currently starting "Phase 1" to make it a little more professional: 

- [x] Internal SSH instead of exec
- [] JSON to DB
- [] API Endpoints


## Setup
`git clone https://github.com/RaynardGerraldo/SimplePAM ; cd SimplePAM ; go mod tidy ; go build`

## Usage

`./SimplePAM admin init` -- initialize admin and server (currently only your localhost)

`./SimplePAM admin add-user <name>` -- add first user for ssh access to server

`./SimplePAM user <user's name>` -- login to allowed server (currently only server-prod)
