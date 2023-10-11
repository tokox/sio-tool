package judge

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Arapak/sio-tool/util"
	"github.com/mitchellh/go-homedir"
)

var options = " --mount-namespace off" +
	" --pid-namespace off" +
	" --uts-namespace off" +
	" --ipc-namespace off" +
	" --net-namespace off" +
	" --capability-drop off --user-namespace off" +
	" -s" +
	" -m 1000000"

const timelimit = time.Second * 10

//go:embed sio2jail
var sio2jail []byte

var sio2jailPath = "~/.st/sio2jail"

const sio2jailCommand = "%v -f 3 --rtimelimit %vms -o oiaug %v -- %v 3> %v"

func InstallSio2Jail() (err error) {
	sio2jailPath, err = homedir.Expand(sio2jailPath)
	if err != nil {
		return
	}
	return os.WriteFile(sio2jailPath, sio2jail, 0755)
}

func RunProcessWithOiejq(processID, command string, input io.Reader) (oiejqProcessInfo *ProcessInfo, err error) {
	sio2jailPath, err = homedir.Expand(sio2jailPath)
	if err != nil {
		return
	}
	if !util.FileExists(sio2jailPath) {
		err = InstallSio2Jail()
		if err != nil {
			return
		}
	}

	oiejqResults, err := os.CreateTemp(os.TempDir(), "sio2jail-")
	if err != nil {
		return
	}
	defer os.Remove(oiejqResults.Name())

	oiejqCommand := fmt.Sprintf(sio2jailCommand, sio2jailPath, timelimit.Milliseconds(), options, command, oiejqResults.Name())
	processInfo, err := RunProcess(processID, oiejqCommand, input, oiejqResults)
	if err != nil {
		return
	}
	result, err := os.ReadFile(oiejqResults.Name())
	if err != nil {
		return
	}
	results := strings.Split(string(result), " ")
	timeMiliseconds, err := strconv.ParseFloat(results[2], 64)
	if err != nil {
		return
	}
	memory, err := strconv.Atoi(results[4])
	if err != nil {
		return
	}

	return &ProcessInfo{timeMiliseconds / 1000, parseMemory(uint64(memory) * 1024), processInfo.Output}, nil
}
