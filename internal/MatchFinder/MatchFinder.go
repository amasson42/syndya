package MatchFinder

import (
	"log"
	"os"
	"syndya/pkg/Models"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type MatchupFunction func(ids []int)

type MatchFinder struct {
	playersBank     Models.SearchingPlayersBank
	scriptPath      string
	luaState        *lua.LState
	resetEachLoop   bool
	MatchupDelegate MatchupFunction
}

func NewMatchFinder(
	playersBank Models.SearchingPlayersBank,
	scriptPath string,
	resetEachLoop bool,
) (*MatchFinder, error) {
	mf := &MatchFinder{
		playersBank:   playersBank,
		scriptPath:    scriptPath,
		luaState:      nil,
		resetEachLoop: resetEachLoop,
	}

	err := mf.reloadScript()
	if err != nil {
		return nil, err
	}

	return mf, nil
}

func (mf *MatchFinder) AsyncRunLoop(timeInterval int) {
	go func() {
		for {
			time.Sleep(time.Duration(timeInterval) * time.Millisecond)
			mf.RunOnce()

			if mf.resetEachLoop {
				mf.reloadScript()
			}
		}
	}()
}

func (mf *MatchFinder) RunOnce() {
	L := mf.luaState
	if L == nil {
		return
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("start"),
		NRet:    0,
		Protect: true,
	}); err != nil {
		log.Printf("[LUA]: %v\n", err)
	}

	mf.playersBank.ForEach(func(player *Models.SearchingPlayer) {

		playerUserData := L.NewUserData()
		playerUserData.Value = player
		L.Push(playerUserData)

		if err := L.CallByParam(lua.P{
			Fn:      L.GetGlobal("__process"),
			NRet:    0,
			Protect: true,
		}, playerUserData); err != nil {
			log.Printf("[LUA]: %v\n", err)
		}

	})

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("finish"),
		NRet:    0,
		Protect: true,
	}); err != nil {
		log.Printf("[LUA]: %v\n", err)
	}
}

func (mf *MatchFinder) reloadScript() error {
	if mf.luaState != nil {
		mf.luaState.Close()
		mf.luaState = nil
	}

	L := lua.NewState()

	L.SetGlobal("__cast_go_player", L.NewFunction(func(L *lua.LState) int {
		player := L.CheckUserData(1).Value.(*Models.SearchingPlayer)

		playerTable := L.NewTable()

		playerTable.RawSetString("searchId", lua.LNumber(player.ID))
		playerTable.RawSetString("waitTime", lua.LNumber(time.Now().Unix()-player.TimeStamp))

		metaTable := L.NewTable()
		for k, v := range player.MetaData {
			metaTable.RawSetString(k, lua.LString(v))
		}
		playerTable.RawSetString("metaDatas", metaTable)

		L.Push(playerTable)

		return 1
	}))

	scriptContent, err := os.ReadFile(mf.scriptPath)
	if err != nil {
		return err
	}

	if err := L.DoString(string(scriptContent)); err != nil {
		return err
	}

	if err := L.DoString(`
	function __process(goPlayer)
		local player = __cast_go_player(goPlayer)
		process(player)
	end
	`); err != nil {
		return err
	}

	L.SetGlobal("matchup", L.NewFunction(func(L *lua.LState) int {
		nArgs := L.GetTop()
		if nArgs != 1 {
			log.Println("Error: matchup expects 1 argument")
			return 0
		}

		if table := L.Get(1); table.Type() == lua.LTTable {
			playersIds := []int{}
			table.(*lua.LTable).ForEach(func(key lua.LValue, value lua.LValue) {
				if num, ok := value.(lua.LNumber); ok {
					playersIds = append(playersIds, int(num))
				}
			})
			if mf.MatchupDelegate != nil {
				mf.MatchupDelegate(playersIds)
			}
		} else {
			log.Printf("Nope %v\n", L.Get(1).Type())
		}
		return 0
	}))

	mf.luaState = L
	return nil
}
