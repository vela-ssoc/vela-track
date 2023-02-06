package track

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

//var cmm = command{
//	exe:  "share\\win\\tools\\handle.exe",
//	hash: "b7976cfa4763e744cbea8ec4e462e185",
//}

var re = regexp.MustCompile(`\s+[a-zA-Z0-9]+\:\s+(\w*)\s+(.*)\s$`)

type command struct {
	exe  string
	hash string
}

func newSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func trim(line string, sec *Section) bool {
	m := re.FindStringSubmatch(line)
	if len(m) < 3 {
		return false
	}

	t := m[1]
	v := m[2]

	i := strings.Index(v, "  ")
	if i < 0 {
		sec.Typ, sec.Value = strings.ToLower(t), v
		return true
	}

	sec.Typ, sec.Value = strings.ToLower(t), v[i+3:]
	return true
}

func HandleExecutable() *command {
	info, err := xEnv.Third("handle.exe")
	if err != nil {
		return &command{}
	}

	return &command{
		exe:  info.Path(),
		hash: info.Hash,
	}
}

func (tk *track) verbose(out io.ReadCloser) {
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		sec := Section{}
		if !trim(line, &sec) {
			goto next
		} else {
			sec.Pid = tk.pid
			sec.Exe = tk.exe
			sec.Args = tk.args
		}

		tk.incr()
		tk.append(sec)

	next:
		if er := scanner.Err(); er != nil {
			xEnv.Errorf("track handle.exe scan fail %v", er)
			return
		}
	}
}

func (tk *track) Command(arg ...string) *exec.Cmd {
	cmm := HandleExecutable()
	if cmm.exe == "" {
		return nil
	}

	cmd := exec.Command(cmm.exe, arg...)
	cmd.SysProcAttr = newSysProcAttr()
	return cmd
}

func (tk *track) forkExec(cmd *exec.Cmd, handle func(io.ReadCloser)) {
	if cmd == nil {
		return
	}

	defer func() {
		if cmd.Process == nil {
			return
		}
		cmd.Process.Kill()
	}()

	out, err := cmd.StdoutPipe()
	if err != nil {
		tk.cause.Try("track command out pipe", err)
		return
	}

	if e := cmd.Start(); e != nil {
		tk.cause.Try("track command", err)
		return
	}

	handle(out)
	cmd.Wait()
}

func (tk *track) Pid() {
	cmm := HandleExecutable()
	if cmm.exe == "" {
		return
	}

	cmd := tk.Command("/accepteula", "-p", tk.pid2str(), "-nobanner")
	tk.forkExec(cmd, tk.verbose)
}
