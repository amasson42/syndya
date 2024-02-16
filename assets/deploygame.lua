
function deploy(players)
    for i, player in pairs(players) do
        print(i .. ":")
        print("  id=" .. player.id)
        print("  rating=" .. player.metaDatas.rating)
    end
    os.execute("sleep 20")
    return "gameip:gameport"
end
