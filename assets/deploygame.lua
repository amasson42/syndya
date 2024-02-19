
function deploy(players)

    local gamedeployerService = getenv("DEPLOYER_SERVICE")
    
    httpRequest("POST", gamedeployerService, {[""]})
    for i, player in pairs(players) do
        print(i .. ":")
        print("  id=" .. player.id)
        print("  rating=" .. player.metaDatas.rating)
    end
    os.execute("sleep 20")
    return "127.0.0.1:7777"
end
