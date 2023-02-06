//go:build linux || plan9 || freebsd || solaris
// +build linux plan9 freebsd solaris

package track

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (tk *track) filepath(pid int32) string {
	return filepath.Join("/proc", strconv.Itoa(int(pid)), "fd")
}

func (tk *track) readdir(pid int32) ([]string, error) {
	path := tk.filepath(pid)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open dir %s", path)
	}
	defer f.Close()

	return f.Readdirnames(-1)
}

func (tk *track) readlink(dir, fd string) (string, error) {
	return os.Readlink(filepath.Join(dir, fd))
}

func (tk *track) Pid() ([]Section, error) {
	sym, err := tk.readdir(tk.pid)
	if err != nil {
		return nil, err
	}

	var data []Section
	dir := filepath.Join("/proc", tk.pid2str(), "fd")
	tk.total = int32(len(sym))

	for _, fd := range sym {
		id := auxlib.ToInt(fd)
		p, e := tk.readlink(dir, fd)
		if e != nil {
			continue
		}
		switch {
		case p == "/dev/null":
			continue
		case strings.HasPrefix(p, "anon_inode:"):
			//tk.append(Section{"inode", p[12 : len(p)-1]})
			continue

		case strings.HasPrefix(p, "socket:"):
			if tk.all || tk.socket {
				tk.append(Section{
					ID:    id,
					Pid:   tk.pid,
					Exe:   tk.exe,
					Typ:   "socket",
					Args:  tk.args,
					Value: p,
				})
			}
			continue

		case strings.HasPrefix(p, "pipe:"):
			if tk.all || tk.socket {
				tk.append(Section{
					ID:    id,
					Pid:   tk.pid,
					Exe:   tk.exe,
					Typ:   "pipe",
					Args:  tk.args,
					Value: p,
				})
			}
			continue

		default:
			tk.append(Section{
				ID:    id,
				Pid:   tk.pid,
				Exe:   tk.exe,
				Typ:   "file",
				Args:  tk.args,
				Value: p,
			})
		}
	}

	return data, nil
}
