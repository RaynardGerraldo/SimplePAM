## SimplePAM

a really basic implementation of a Privileged Access Management 

## Setup
`git clone https://github.com/RaynardGerraldo/SimplePAM ; cd SimplePAM ; go mod tidy ; go build`

## Usage

Create secrets/jwt_secret.txt, with your random jwt secret. Then:

`go run api/*` -- run this on a seperate terminal

`./SimplePAM admin init` -- initialize admin and server (currently only your localhost)

`./SimplePAM admin add-user <name>` -- add first user for ssh access to server

`./SimplePAM admin add-server` -- add new server

`./SimplePAM admin srv-to-user` -- assign server to user

http://127.0.0.1:8080/web/ -- PAM frontend
