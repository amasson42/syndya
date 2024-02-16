package GameDeployer

import (
	"errors"
	"syndya/internal/LuaExtensionSyndya"
	"syndya/pkg/LuaExtension"
	"syndya/pkg/Models"

	lua "github.com/yuin/gopher-lua"
)

type GameDeployer struct {
	scriptPath string
}

func NewGameDeployer(
	scriptPath string,
) *GameDeployer {
	gd := &GameDeployer{
		scriptPath: scriptPath,
	}

	return gd
}

// reloadScript reloads the Lua script.
func (gd *GameDeployer) LoadScript() (*lua.LState, error) {
	L := lua.NewState()

	LuaExtension.AddGetenvFunction(L)
	LuaExtension.AddHttpRequestFunction(L)
	LuaExtension.AddJsonFunction(L)

	LuaExtensionSyndya.AddCastGoPlayerFunction(L)

	if err := L.DoFile(gd.scriptPath); err != nil {
		return nil, err
	}

	if err := L.DoString(`
	function __deploy(goPlayers)
		local players = {}

		for i, goPlayer in ipairs(goPlayers) do
			local player = __cast_go_player(goPlayer)

			table.insert(players, player)
		end

		return deploy(players)
	end
	`); err != nil {
		return nil, err
	}

	return L, nil
}

func (gd *GameDeployer) DeployGame(players []*Models.SearchingPlayer) (*string, error) {
	L, err := gd.LoadScript()
	if err != nil {
		return nil, err
	}
	defer L.Close()

	playerTable := L.CreateTable(len(players), 0)
	for i, player := range players {
		playerUserData := L.NewUserData()
		playerUserData.Value = player
		playerTable.RawSetInt(i+1, playerUserData)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("__deploy"),
		NRet:    1,
		Protect: true,
	}, playerTable); err != nil {
		return nil, err
	}

	result := L.Get(-1)
	L.Pop(1)

	if str, ok := result.(lua.LString); ok {
		gameaddr := str.String()
		return &gameaddr, nil
	} else {
		return nil, errors.New("result is not a string")
	}
}
