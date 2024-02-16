
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

| Name | Type | Description |
| ---- | ---- | ---- |
| `LISTEN_HOST` | string | Http server listen host |
| `LISTEN_PORT` | integer | Http server listen port |
| `METADATAS_LIST` | string | Comma separated list of requested metadata for the algorithm to function |
| `METADATAS_REQUIRED` | bool | Option to ignore players with missing metadatas when processing matchmakings |
| `METADATAS_REVIVEPERIOD` | integer | Time in milliseconds between each request from the server to fetch missing metadatas |
| `MATCHFINDER_LUASCRIPT` | string | Path to the lua script file that manage matchmaking (details at [[Matchmaking script]]). If none is given, this service won't run. |
| `MATCHFINDER_TIMEINTERVAL` | integer | Time in milliseconds between each call to the custom script with all searching players |
| `MATCHFINDER_RESETSTATE` | bool | Fully destroy and re create the Lua state between each matchmaking loops |
| `GAMEDEPLOYER_LUASCRIPT` | string | Path to the lua script file that manage game deployement (details at [[Game Deploying script]]). If none is given, this service won't run. |
