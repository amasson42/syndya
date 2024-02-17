The service employs algorithms that may vary based on game or infrastructure requirements. As a result, certain functionalities are implemented using external scripts.

There is an built-in interpreter for the Lua scripting language to accommodate these external scripts. Lua files with extended functionalities and specific behaviors are utilized to enhance the service's capabilities.

> [lua](https://www.lua.org/) is a famous scripting language used in many areas including video game plugins.

To enable the service to utilize these scripts effectively, certain environment variables must be configured to define their respective paths.

## Matchfinder Script

The Matchfinder script is used to manage matchmaking operations within the service.

### Configuration

- **Environment Variable**: `MATCHFINDER_LUASCRIPT`
- **Time Interval**: The default time interval for executing the script is 5000 milliseconds. This can be modified with `MATCHFINDER_TIMEINTERVAL`

### Execution Steps

1. The Lua function `start()` is executed at the beginning of each iteration.
2. The service iterates through all searching players, creating a Lua table for each player's data, and executes the Lua function `process(player)`.
3. The Lua function `finish()` is executed at the end of each iteration.

The execution is single threaded and the time interval for the next execution is only triggered after the first execution is finished.

### Player Structure

The `process(player)` function receives a player object with the following fields:

- `id`: (integer) Identification number of the searching player.
- `waitTime`: (integer) Amount of seconds the player has been waiting.
- `metaDatas`: (table) Table of key-value pairs representing [[Players Metadatas]].

### Matching Up

After processing a player, if a potential match is found, the service can proceed to start a game. The `matchup(table)` function is available for this purpose, which takes a list of player IDs as input.

### Script Example

```lua
playerIds = {}

-- Executed once at the beginning of the iteration
function start()
  playerIds = {}
  print("Starting iteration")
end

-- Executed for each player
function process(player)
  print("Processing player " .. player.id)
  print("Waiting since " .. player.waitTime .. " seconds")
  print("Metadata rating: " .. player.metaDatas.rating)
  table.insert(playerIds, player.id)

  if #playerIds >= 4 then
    matchup(playerIds)
    playerIds = {}
  end
end

-- Executed once at the end of the iteration
function finish()
  print("End iteration")
end
```

## Game deployer script

The Game Deployer script manages the deployment of games within the service.

### Configuration

- **Environment Variable**: `GAMEDEPLOYER_LUASCRIPT`



# Extended lua functions

Additional functions are added to scripts to allow for more possibilities or simplify processes.

These functions leverage compiled Go features for improved performance.

| Function | Return Type | Description | Example |
| ---- | ---- | ---- | ---- |
| `getenv(string)` | `string` | Fetches the value of an environment variable. | `user = getenv("USER")` |
| `httpRequest(string, string, table, string)` | `string` | Makes an HTTP request and returns the body of the response. | `response = httpRequest("POST", "https://example.com/api", {["Content-Type"] = "application/json"}, '{"key": "value"}')` |
| `json(string)` | `table` | Parses a string into a JSON object. | `jsonObject = json('{"name": "John", "age": 30, "city": "New York"}')`<br> |
