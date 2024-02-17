playerIds = {}

-- Used once at the beginning of the iteration
function start()
    playerIds = {}
    print("Starting iteration")
end

-- Used for each players
function process(player)
    print("Processing player " .. player.id)
    print("waiting since " .. player.waitTime)
    print("metadatas " .. player.metaDatas.rating)
    table.insert(playerIds, player.id)

    if #playerIds >= 4 then
        matchup(playerIds)
        playerIds = {}
    end
end

-- Used once at the end of the iteration
function finish()
    print("End iteration")
end