package track

import (
	cond "github.com/vela-ssoc/vela-cond"
	auxlib2 "github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
)

type LOption struct {
	file   bool
	socket bool
	all    bool
}

func (opt *LOption) String() string                         { return auxlib2.B2S(opt.Byte()) }
func (opt *LOption) Type() lua.LValueType                   { return lua.LTObject }
func (opt *LOption) AssertFloat64() (float64, bool)         { return 0, false }
func (opt *LOption) AssertString() (string, bool)           { return "", false }
func (opt *LOption) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (opt *LOption) Peek() lua.LValue                       { return opt }

func (opt *LOption) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("file", opt.file)
	enc.KV("socket", opt.socket)
	enc.KV("all", opt.all)
	enc.End("}")
	return enc.Bytes()
}

func (opt *LOption) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "socket":
		opt.socket = lua.IsTrue(val)
	case "all":
		opt.all = lua.IsTrue(val)
	case "file":
		opt.file = lua.IsTrue(val)
	}
}

func (opt *LOption) doL(L *lua.LState) int {
	tka := &tracks{cnd: cond.CheckMany(L), opt: opt}
	tka.scan()
	L.Push(tka)
	return 1
}

func (opt *LOption) withPidL(L *lua.LState) int {
	pid := L.IsInt(1)
	cnd := cond.CheckMany(L, cond.Seek(1))
	tk := newTrack(withPid(int32(pid)), withCnd(cnd), withOption(opt))
	tk.lookup()
	tk.Pid()
	return 1
}

func (opt *LOption) withNameL(L *lua.LState) int {
	L.Push(newTrackNameByOption(L, opt))
	return 1
}

func (opt *LOption) withKwL(L *lua.LState) int {
	L.Push(newTracksKeyWoldByOption(L, opt))
	return 1
}

func (opt *LOption) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "with":
		return lua.NewFunction(opt.doL)
	case "pid":
		return lua.NewFunction(opt.withPidL)
	case "name":
		return lua.NewFunction(opt.withNameL)
	case "kw":
		return lua.NewFunction(opt.withKwL)
	}

	return lua.LNil
}

func newLuaOption(L *lua.LState) *LOption {
	opt := &LOption{}

	L.Callback(func(lv lua.LValue) (stop bool) {
		if lv.Type() != lua.LTString {
			return true
		}
		k, v := auxlib2.ParamLValue(lv.String())
		opt.NewIndex(L, k, v)
		return false
	})

	return opt
}
