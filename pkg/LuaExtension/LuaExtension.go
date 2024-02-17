package LuaExtension

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func AddGetenvFunction(L *lua.LState) {
	L.SetGlobal("getenv", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		value := os.Getenv(key)
		L.Push(lua.LString(value))
		return 1
	}))
}

func AddHttpRequestFunction(L *lua.LState) {
	L.SetGlobal("httpRequest", L.NewFunction(func(L *lua.LState) int {
		verb := L.ToString(1)
		url := L.ToString(2)

		req, err := http.NewRequest(verb, url, nil)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Error creating HTTP request: %v", err)))
			return 2
		}

		if L.GetTop() >= 3 {
			headersTable := L.ToTable(3)
			headers := make(http.Header)
			headersTable.ForEach(func(key, value lua.LValue) {
				headers.Add(lua.LVAsString(key), lua.LVAsString(value))
			})
			req.Header = headers
		}

		if L.GetTop() >= 4 {
			bodyStr := L.ToString(4)
			req.Body = io.NopCloser(strings.NewReader(bodyStr))
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Error making HTTP request: %v", err)))
			return 2
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Error reading response body: %v", err)))
			return 2
		}

		L.Push(lua.LString(body))

		return 1
	}))
}

func AddJsonFunction(L *lua.LState) {
	L.SetGlobal("json", L.NewFunction(func(L *lua.LState) int {
		jsonStr := L.CheckString(1)

		var jsonObj interface{}
		err := json.Unmarshal([]byte(jsonStr), &jsonObj)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("Error parsing JSON: %v", err)))
			return 2
		}

		luaValue := convertToLuaValue(L, jsonObj)

		L.Push(luaValue)
		return 1

	}))
}

func convertToLuaValue(L *lua.LState, val interface{}) lua.LValue {
	switch v := val.(type) {
	case map[string]interface{}:
		luaTable := L.NewTable()
		for key, subVal := range v {
			L.SetField(luaTable, key, convertToLuaValue(L, subVal))
		}
		return luaTable
	case []interface{}:
		luaTable := L.NewTable()
		for i, subVal := range v {
			L.SetField(luaTable, fmt.Sprintf("[%d]", i+1), convertToLuaValue(L, subVal))
		}
		return luaTable
	case string:
		return lua.LString(v)
	case float64:
		return lua.LNumber(v)
	case bool:
		return lua.LBool(v)
	default:
		return lua.LNil
	}
}
