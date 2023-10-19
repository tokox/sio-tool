package judge

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Arapak/sio-tool/util"
	"github.com/mitchellh/go-homedir"
)

type OiejqOptions struct {
	MemorylimitInMegaBytes string
	TimeLimitInSeconds     string
}

var options = " --mount-namespace off" +
	" --pid-namespace off" +
	" --uts-namespace off" +
	" --ipc-namespace off" +
	" --net-namespace off" +
	" --capability-drop off --user-namespace off" +
	" -s"

const defaultTimeLimit = "10"
const defaultMemoryLimitInMegaBytes = "1024"

//go:embed sio2jail
var sio2jail []byte

var sio2jailPath = "~/.st/sio2jail"

const sio2jailCommand = "%v -f 3 --instruction-count-limit %vg -o oiaug %v --memory-limit %vM -- %v 3> %v"

func InstallSio2Jail() (err error) {
	sio2jailPath, err = homedir.Expand(sio2jailPath)
	if err != nil {
		return
	}
	return os.WriteFile(sio2jailPath, sio2jail, 0755)
}

const ErrorInvalidOiejqResults = "invalid oiejq results returned"

func readOiejqOutput(processID, path string) (*ProcessInfo, error) {
	result, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(result), "\n")
	if len(lines) != 3 {
		return nil, fmt.Errorf(ErrorInvalidOiejqResults)
	}
	results := strings.Split(lines[0], " ")
	message := lines[1]
	if results[0] == "RE" || results[0] == "RV" {
		return nil, fmt.Errorf("runtime error #%v ... %v", processID, message)
	} else if results[0] == "TLE" {
		return nil, fmt.Errorf("time limit exceeded #%v", processID)
	} else if results[0] == "MLE" {
		return nil, fmt.Errorf("memory limit exceeded #%v", processID)
	} else if results[0] == "OLE" {
		return nil, fmt.Errorf("output limit exceeded #%v", processID)
	} else if results[0] != "OK" {
		return nil, fmt.Errorf("invalid oiejq status returned")
	}

	timeMiliseconds, err := strconv.ParseFloat(results[2], 64)
	if err != nil {
		return nil, err
	}
	memory, err := strconv.Atoi(results[4])
	if err != nil {
		return nil, err
	}

	return &ProcessInfo{timeMiliseconds / 1000, parseMemory(uint64(memory) * 1024), []byte{}}, nil
}

func RunProcessWithOiejq(processID, command string, input io.Reader, oiejqOptions *OiejqOptions) (oiejqProcessInfo *ProcessInfo, err error) {
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

	if oiejqOptions.MemorylimitInMegaBytes == "" {
		oiejqOptions.MemorylimitInMegaBytes = defaultMemoryLimitInMegaBytes
	}
	if oiejqOptions.TimeLimitInSeconds == "" {
		oiejqOptions.TimeLimitInSeconds = defaultTimeLimit
	}

	oiejqCommand := fmt.Sprintf(sio2jailCommand, sio2jailPath, oiejqOptions.TimeLimitInSeconds, options, oiejqOptions.MemorylimitInMegaBytes, command, oiejqResults.Name())
	processInfo, processErr := RunProcess(processID, oiejqCommand, input, oiejqResults)
	oiejqProcessInfo, err = readOiejqOutput(processID, oiejqResults.Name())
	if err != nil {
		if err.Error() == ErrorInvalidOiejqResults {
			if processErr != nil {
				return nil, processErr
			} else if processInfo == nil {
				return nil, fmt.Errorf("runtime error #%v", processID)
			} else {
				return nil, errors.New(string(processInfo.Output))
			}
		}
		return
	}
	if processErr != nil {
		return nil, processErr
	}
	oiejqProcessInfo.Output = processInfo.Output
	return
}
