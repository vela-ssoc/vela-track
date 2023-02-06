package track

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

func newLuaTrackPidL(L *lua.LState) int {
	pid := L.IsInt(1)
	cnd := cond.CheckMany(L, cond.Seek(1))
	L.Push(newTrackByPid(int32(pid), cnd))
	return 1
}

func newLuaTrackAllL(L *lua.LState) int {
	tka := newLuaTrackALL(L)
	L.Push(tka)
	return 1
}

func newLuaTrackKwL(L *lua.LState) int {
	L.Push(newTracksKeyWold(L))
	return 1
}

func newLuaTrackNameL(L *lua.LState) int {
	L.Push(newTrackName(L))
	return 1
}

func newLuaTrackByParamL(L *lua.LState) int {
	L.Push(newLuaOption(L))
	return 1
}

/*
	local v = vela.track.pid()
	local v = vela.track.p("socket:true" , "all:true" , "file:true").do("").pipe()
	vela.track.p("socket:true" , "all:true" , "file:true").pid(1 , "").pipe()
	vela.track.p("socket:true" , "all:true" , "file:true").name("" , "").pipe()
	vela.track.p("socket:true" , "all:true" , "file:true").kw("" , "").pipe()
*/

func WithEnv(env vela.Environment) {
	xEnv = env
	kv := lua.NewUserKV()
	kv.Set("pid", lua.NewFunction(newLuaTrackPidL))
	kv.Set("kw", lua.NewFunction(newLuaTrackKwL))
	kv.Set("name", lua.NewFunction(newLuaTrackNameL))
	kv.Set("p", lua.NewFunction(newLuaTrackByParamL))

	xEnv.Set("track",
		lua.NewExport("vela.track.export",
			lua.WithTable(kv),
			lua.WithFunc(newLuaTrackAllL)))
}
