package judge

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/Arapak/sio-tool/util"
	"github.com/shirou/gopsutil/process"
)

func Plain(raw []byte) string {
	buf := bufio.NewScanner(bytes.NewReader(raw))
	var b bytes.Buffer
	newline := []byte{'\n'}
	for buf.Scan() {
		line := bytes.TrimSpace(buf.Bytes())
		if len(line) != 0 {
			b.Write(line)
			b.Write(newline)
		}
	}
	return b.String()
}

type ProcessInfo struct {
	Status            VerdictStatus
	TimeInSeconds     float64
	MemoryInMegabytes float64
	Output            []byte
}

func RunProcess(command string, input io.Reader, extrafile *os.File) (ProcessInfo, error) {
	var o bytes.Buffer
	output := io.Writer(&o)
	var e bytes.Buffer
	stderr := io.Writer(&e)

	cmds := util.SplitCmd(command)

	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = stderr
	if extrafile != nil {
		cmd.ExtraFiles = append(cmd.ExtraFiles, extrafile)
	}
	if err := cmd.Start(); err != nil {
		return ProcessInfo{RE, 0, 0, []byte{}}, err
	}

	pid := int32(cmd.Process.Pid)
	maxMemory := uint64(0)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			if err != nil {
				return ProcessInfo{RE, 0, float64(maxMemory) / (1024.0 * 1024.0), []byte{}}, err
			}
			running = false
		default:
			p, err := process.NewProcess(pid)
			if err == nil {
				m, err := p.MemoryInfo()
				if err == nil && m.RSS > maxMemory {
					maxMemory = m.RSS
				}
			}
			if extrafile != nil && bytes.Contains(e.Bytes(), []byte("Exception occurred: System error occured: perf event open failed: Permission denied: error 13: Permission denied")) {
				cmd.Process.Kill()
				return ProcessInfo{PERF, 0, float64(maxMemory) / (1024.0 * 1024.0), []byte{}}, nil
			}
		}
	}
	return ProcessInfo{OK, cmd.ProcessState.UserTime().Seconds(), float64(maxMemory) / (1024.0 * 1024.0), o.Bytes()}, nil
}
