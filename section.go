package track

import (
	"bytes"
	"fmt"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"path/filepath"
	"strconv"
	"strings"
)

type Section struct {
	ID    int
	Pid   int32
	User  string
	Exe   string
	Typ   string
	Value string
	Args  string
}

func (s *Section) String() string                         { return lua.B2S(s.Byte()) }
func (s *Section) Type() lua.LValueType                   { return lua.LTObject }
func (s *Section) AssertFloat64() (float64, bool)         { return 0, false }
func (s *Section) AssertString() (string, bool)           { return "", false }
func (s *Section) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (s *Section) Peek() lua.LValue                       { return s }

func (s *Section) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("id", s.ID)
	enc.KV("pid", s.Pid)
	enc.KV("type", s.Typ)
	enc.KV("value", s.Value)
	enc.KV("name", s.Name())
	enc.KV("exe", s.Exe)
	enc.End("}")
	return enc.Bytes()
}

func (s *Section) Name() string {
	if s.Exe == "" {
		return ""
	}
	return filepath.Base(s.Exe)
}

func (s *Section) Raw() string {
	var buf bytes.Buffer
	buf.WriteString(strconv.Itoa(int(s.Pid)))
	buf.WriteByte(' ')
	buf.WriteString(s.Value)
	buf.WriteByte(' ')
	buf.WriteString(s.User)
	buf.WriteByte(' ')
	buf.WriteString(s.Exe)
	buf.WriteByte(' ')
	buf.WriteString(s.Typ)
	return buf.String()
}

func (s *Section) inode() string {
	switch s.Typ {
	case "socket", "pipe":

		idx := strings.Index(s.Value, ":[")
		if idx == -1 {
			return ""
		}

		n := len(s.Value)
		return s.Value[idx+2 : n-1]

	default:
		return ""
	}
}

func (s *Section) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "id":
		return lua.LInt(s.ID)
	case "type":
		return lua.S2L(s.Typ)
	case "value":
		return lua.S2L(s.Value)
	case "inode":
		return lua.S2L(s.inode())
	case "pid":
		return lua.LNumber(s.Pid)
	case "exe":
		return lua.LString(s.Exe)
	case "name":
		return lua.LString(s.Name())
	case "ext":
		return lua.LString(filepath.Ext(s.Exe))
	case "info":
		return lua.S2L(fmt.Sprintf("pid:%d name:%s type:%s value:%s", s.Pid, s.Name(), s.Typ, s.Value))
	case "raw":
		return lua.LString(s.Raw())
	}

	return lua.LNil
}
