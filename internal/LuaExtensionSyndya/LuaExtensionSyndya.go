package LuaExtensionSyndya

import (
	"syndya/pkg/Models"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func AddCastGoPlayerFunction(L *lua.LState) {
	L.SetGlobal("__cast_go_player", L.NewFunction(func(L *lua.LState) int {
		player := L.CheckUserData(1).Value.(*Models.SearchingPlayer)

		playerTable := L.NewTable()

		playerTable.RawSetString("id", lua.LNumber(player.ID))
		playerTable.RawSetString("waitTime", lua.LNumber(time.Now().Unix()-player.TimeStamp))

		metaTable := L.NewTable()
		for k, v := range player.MetaData {
			metaTable.RawSetString(k, lua.LString(v))
		}
		playerTable.RawSetString("metaDatas", metaTable)

		L.Push(playerTable)

		return 1
	}))
}
