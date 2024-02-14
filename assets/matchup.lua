
players = {}

-- Used once at the beginning of the iteration
function start()
    players = {}
    print("Starting iteration")
end

-- Used for each players
function process(player)
    print("Processing player " .. player.searchId)
    print("waiting since " .. player.waitTime)
    print("metadatas " .. player.metaDatas.rating)
    table.insert(players, player.searchId)

    if #players >= 4 then
        matchup(players)
        players = {}
    end
end

-- Used once at the end of the iteration
function finish()
    print("End iteration")
end
