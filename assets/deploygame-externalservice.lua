
-- This bellow example is an example of how to deploy a game using an external service.
-- We expect the user to send the address of this external service in the DEPLOYER_SERVICE environment variable
-- The service is expected to start a game when request and return in the response body a json with a field 'ipaddr' of the deployed game
function deploy(players)

    -- Fetch the a game deployer microservice
    local serviceUrl = getenv("DEPLOYER_SERVICE")

    -- Post a request to deploy a game and expect a response
    local result = httpRequest(
        "POST",
        serviceUrl,
        {
            ["Content-Type"] = "application/json"
        },
        string.format('{"playersCount": "%d"}', #players)
    )

    -- Parse the response JSON formatted into a readable table
    local resultJson = json(result)

    -- Extract the address of the game from the json
    return resultJson.ipaddr

end
