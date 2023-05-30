## 👋 Intro
ALO is simple low lever demo of a tcp connection with realtime communication. It was made for an assigment a long time ago - hence the incompletness

## Requirements

 - go1.20 (mandatory)
 - Docker (optional)


## Running the project
### Dependencies
- (optional) Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- (recommended) Install [go 1.14 and up](https://golang.org/doc/install)
### Commands

1. go get ./...
2. make start.server
3. _on another console_, telnet localhost 8080 <- and type the payload to test the server

### Commands with docker

1. make docker.start
2. on another console, telnet localhost 8080 <- and type the payload to test the server

## Tests

The project contains unit tests (partial), to run the tests

`make tests`

## Available make commands

| command            | description                                                              |
|--------------------|--------------------------------------------------------------------------|
| make start.server  | start the game server - Environment: $VERBOSE $SERVER_PORT               |
| make start.client  | (not implemented yet) start the client to test the game server           |
| make tools         | install developer tools                                                  |
| make build         | only build the client and the server (binaries installed in `.bin`)      |
| make assets        | generate marshalling code for the models                                 |
| make tests         | run the unit tests                                                       |
| make lint          | lint the whole project                                                   |
| make docker.start  | run the application inside containers                                    |
| make docker.stop   | stop the docker stack                                                    |
| make docker.build  | only build the containers                                                |
| make docker.logs   | attach and follow containers log                                         |
| make docker.client | (not implement yet) start and attach the current console with the client |
| make docker.lint   | lint the whole project using docker instead                              |

## Project Structure

| name                   | description                                                                           |
|------------------------|---------------------------------------------------------------------------------------|
| .bin                   | holds generated binaries (client, server)                                             |
| cmd                    | contains the server and client cli code                                               |
| internal/mock          | generated mocks to help with testing                                                  |
| pkg/model              | models used to communicate between servers and clients                                |
| pkg/server             | logic necessary to for the game server                                                |
| pkg/store              | abstraction directing the use and implementation of any store working with the server |
| pkg/store/light        | very lightweight in memory store built for the game                                   |
| pkg/store/{redis;psql} | for example purpose (no code) 

## Server help

```
$ ./alo -h
NAME:
   ALO

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, -v           verbose (default: false) [$VERBOSE]
   --port value, -p value  set the port to listen to (default: 8080) [$SERVER_PORT]
   --help, -h              show help (default: false)
```
## Requirements
1. Start listening for tcp/udp connections.
2. Be able to accept connections.
3. Read json payload `{"user_id": 1, "friends": [2, 3, 4]}`
3. After establishing successful connection - "store" it in memory the way you like.
4. When another connection established with the `user_id` from the list of any other user's "friends" section, they should be notified about it with message {"online": true}
5. When the user goes offline, his "friends" (if it has any and any of them online) should receive a message `{"online": false}`

## Dev environment
### Requirement
- go1.20
- macOS - or Docker (wip) - or install manually the tools inside `make tools`

### Setup the environment

1. `make tools`


### Lint

console: `make lint`

docker: `make docker.lint`

## Technical Debt

- Client - interactive client to help to debug the server and manually test it with [go-prompt](https://github.com/c-bata/go-prompt)
- Server benchmark - and race testing
- UDP server implementations
- Unit tests
- Integration tests


## Design decisions

Given the limited time I have I followed a simple and straightforward principle where the server redirects the connections to the router
which will be responsible to decide what to do with connection events. I coded a simple store to ease the friend relationships as all in memory lib I could find have 
very limited queries option. I also made sure to think about performance, and I only used library specifically designed with low allocation and performance in mind:
- github.com/francoispqt/gojay - to be able to decode json on the wire
- github.com/francoispqt/onelog - the fastest json logger 
- github.com/rs/xid - one of the fastest unique ID generator

## Improvements

If I had more time in front of me, I would have designed the server to follow the actors design pattern to ensure better scalability and to avoid any race conditions.
I would provide more interaction with other Storage.

At the moment, the server cannot have multiple instances running at the same time, however we could easily think of possible
implementation in something like a kubernetes cluster once a shared storage has been implemented.

We can also think of possible Master-slave architecture where the master get metrics of its slaves and redirect the connection to them (based on different factors)

I would also consider any RPC abstraction if we don't need UDP, [gRPC](https://grpc.io/) could be a good candidate or [flatbuffer](https://google.github.io/flatbuffers/) 
if we want to be lighter and faster. jsonRPC can be also considered if the platform cannot handle generated proto or flat files