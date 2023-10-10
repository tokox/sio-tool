package judge

import (
	"bufio"
	"bytes"
	"fmt"
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
		b.Write(bytes.TrimSpace(buf.Bytes()))
		b.Write(newline)
	}
	return b.String()
}

func parseMemory(memory uint64) string {
	if memory > 1024*1024 {
		return fmt.Sprintf("%.3fMB", float64(memory)/1024.0/1024.0)
	} else if memory > 1024 {
		return fmt.Sprintf("%.3fKB", float64(memory)/1024.0)
	}
	return fmt.Sprintf("%vB", memory)
}

type ProcessInfo struct {
	Time   float64
	Memory string
	Output []byte
}

func RunProcess(processID, command string, input io.Reader, extrafile *os.File) (*ProcessInfo, error) {
	var o bytes.Buffer
	output := io.Writer(&o)

	cmds := util.SplitCmd(command)

	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	if extrafile != nil {
		cmd.ExtraFiles = append(cmd.ExtraFiles, extrafile)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("runtime error #%v ... %v", processID, err.Error())
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
				return nil, fmt.Errorf("runtime error #%v ... %v", processID, err.Error())
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
		}
	}
	return &ProcessInfo{cmd.ProcessState.UserTime().Seconds(), parseMemory(maxMemory), o.Bytes()}, nil
}
