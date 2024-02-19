# Syndya - Matchmaking Service

This project is a solution designed to streamline the process of pairing players for online games.
The service includes features such as WebSocket communication, custom Lua scripting for matchmaking logic and game deployment functionality, and integration with Redis for efficient data management.

## Features
- **WebSocket Communication**: Players can connect via WebSocket to search for a game with matching players.
- **Matchmaking Logic**: Utilizes a custom Lua script to implement matchmaking algorithms tailored to specific game requirements.
- **Game Deployment**: When matching players are found, another Lua script is triggered to deploy a game instance, and the IP address of the game is sent to the connected players.
- **Administrator Interface**: Endpoints for administrators to monitor connected players and their status.

## Getting Started

An included Makefile contains most build commands

```sh
$ make
help:                    This help dialog.
install:                 Locally fetch dependencies
build:                   Build the local package
build_release:           Build the local package in release mode
run:                     Run the local package
clean:                   Clean the build object and stop dependencies
fclean: clean            Reset the state of the package
kill:                    Stop the program listening on target port
test:                    Execute the unit tests
build_docker:            Build the docker image
test_docker:             Run the test in a docker image
run_docker:              Run the docker image and the dependency
```

The service is written in Golang. Make sure to install it before building it.
> https://go.dev/doc/install

### Build locally

```sh
make build
```

## Connect Players:

Players can connect to the WebSocket endpoint /search to provide their metadata and participate in matchmaking.

## Monitor Players (Optional):

Administrators can monitor connected players and their status by accessing the /players endpoint.

# Contributing
Contributions are welcome! If you have suggestions for improvements or new features, please open an issue or submit a pull request.

# License
This project is licensed under the MIT License - see the LICENSE file for details.
