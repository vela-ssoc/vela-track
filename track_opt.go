package track

import cond "github.com/vela-ssoc/vela-cond"

func withPid(v int32) func(*track) {
	return func(tk *track) {
		tk.pid = v
	}
}

func withCnd(cnd *cond.Cond) func(*track) {
	return func(tk *track) {
		tk.cnd = cnd
	}
}

func withOption(opt *LOption) func(*track) {
	return func(tk *track) {
		if opt == nil {
			return
		}

		tk.socket = opt.socket
		tk.all = opt.all
		tk.file = opt.file
	}

}
