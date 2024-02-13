
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
