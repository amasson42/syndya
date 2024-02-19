---

kanban-plugin: basic

---

## Done

- [ ] hold metadata for players
- [ ] Iterate throught all players to create a match
- [ ] Only process players with full metadata
- [ ] Send game to players
- [ ] Use external lua script to start a game
- [ ] Write documentation on how to write matchmaking script


## Ongoing

- [ ] Documentation on game deployement script


## Pending

- [ ] Remove gin dependency to use standard library
- [ ] Gracefully close connection when player disconnect to join a game
- [ ] Implement redis server usage
- [ ] Communicate with a kubernetes cluster to start a new game
- [ ] use gRPC instead of clocking database read to trigger player game found notification


## Doc

- [ ] Documentation on setting up this service
- [ ] Documentation on scaling up this service
- [ ] Documentation on metadata




%% kanban:settings
```
{"kanban-plugin":"basic"}
```
%%