//go:build linux || plan9 || freebsd || solaris
// +build linux plan9 freebsd solaris

package track

import (
	"fmt"
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/lua"
)

func newTracksKeyWoldByOption(L *lua.LState, opt *LOption) *tracks {
	world := L.CheckString(1)
	cnd := cond.New(fmt.Sprintf("raw cn *%s*", world))
	cnd.CheckMany(L, cond.Seek(1))

	tks := &tracks{cnd: cnd, opt: opt}
	tks.scan()
	return tks
}

func newTracksKeyWold(L *lua.LState) *tracks {
	world := L.CheckString(1)
	cnd := cond.New(fmt.Sprintf("raw cn *%s*", world))
	cnd.CheckMany(L, cond.Seek(1))

	tks := &tracks{cnd: cnd}
	tks.scan()
	return tks
}

func newTrackByName(name string, v *cond.Cond) *tracks {
	cnd := cond.New(fmt.Sprintf("name = %s", name))
	cnd.Merge(v)
	tks := &tracks{cnd: cnd}
	tks.scan()
	return tks
}

func newTrackName(L *lua.LState) *tracks {
	name := L.CheckString(1)
	cnd := cond.New(fmt.Sprintf("name eq %s", name))
	cnd.CheckMany(L, cond.Seek(1))
	tks := &tracks{cnd: cnd}
	tks.scan()
	return tks
}

func newTrackNameByOption(L *lua.LState, opt *LOption) *tracks {
	name := L.CheckString(1)
	cnd := cond.New(fmt.Sprintf("name eq %s", name))
	cnd.CheckMany(L, cond.Seek(1))
	tks := &tracks{cnd: cnd, opt: opt}
	tks.scan()
	return tks
}
