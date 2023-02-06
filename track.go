package track

import (
	"github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/execpt"
	"github.com/vela-ssoc/vela-process"
	"strconv"
	"sync/atomic"
)

type track struct {
	pid    int32
	exe    string
	total  int32
	data   []Section
	socket bool
	file   bool
	all    bool
	cnd    *cond.Cond
	cause  *execpt.Cause
	args   string
}

func (tk *track) pid2str() string {
	return strconv.Itoa(int(tk.pid))
}

func (tk *track) lookup() {

	pro, err := process.Pid(int(tk.pid))
	if err != nil {
		tk.cause.Try("process", err)
		return
	}

	tk.exe = pro.Executable
	tk.args = pro.ArgsToString()
}

func (tk *track) ok() bool {
	return tk.cause.Len() == 0
}

func (tk *track) append(s Section) {
	if tk.cnd.Match(&s) {
		tk.data = append(tk.data, s)
	}
}

func (tk *track) Visit(handle func(s Section)) {
	n := len(tk.data)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		handle(tk.data[i])
	}
}

func (tk *track) Reset() {
	tk.data = nil
}

func (tk *track) incr() {
	atomic.AddInt32(&tk.total, 1)
}

func newTrackByPid(pid int32, cnd *cond.Cond) *track {
	tk := newTrack(withPid(pid), withCnd(cnd))
	tk.lookup()
	tk.Pid()
	return tk
}

func newTrack(opt ...func(*track)) *track {
	ov := &track{cause: execpt.New()}
	for _, fn := range opt {
		fn(ov)
	}
	return ov
}
