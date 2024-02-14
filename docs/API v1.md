
# API Documentation

This document outlines the endpoints and functionalities provided by the API.

## Endpoints

> `GET /search`
- **Description**: Subscribe to look for a game.
- A web socket connection is established

*Server messages*
- `id <playerId>` : Send the player searching id
- `meta <dataKey>` : Kindly ask the user to send the value of the metadata
- `game <gameInfo>` : Send the playable game informations to the player

*Client messages*
- `id` : Kindly request the server to send out searching id
- `meta <string>` : Send to the server a metadata with a key and a value in the format `key=value`

> `GET /players`
- **Description**: Get the list of all searching players

# Environment Variables

`LISTEN_HOST` (String) Http server listen host
`LISTEN_PORT` (Integer) Http server listen port

`METADATAS_LIST` (String) Comma separated list of requested metadata for the algorithm to function
`METADATAS_REVIVEPERIOD` (Integer) Time in milliseconds between each request from the server to fetch missing metadatas

`MATCHFINDER_LUASCRIPT` (String) Path to the lua script file that manage matchmaking (details at [[Matchmaking script]]). If none is given, this service won't run.
`MATCHFINDER_TIMEINTERVAL` (Integer) Time in milliseconds between each call to the custom script with all searching players
`MATCHFINDER_RESETSTATE` (Bool) Fully destroy and re create the lua state between each matchmaking loops
