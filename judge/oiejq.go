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
	if util.FileExists(sio2jailPath) {
		return
	}
	return os.WriteFile(sio2jailPath, sio2jail, 0755)
}

const ErrorInvalidOiejqResults = "invalid oiejq results returned"

func readOiejqOutput(path string) (processInfo ProcessInfo, err error) {
	result, err := os.ReadFile(path)
	if err != nil {
		processInfo.Status = INT
		return
	}
	lines := strings.Split(string(result), "\n")
	if len(lines) != 3 {
		processInfo.Status = INT
		err = errors.New(ErrorInvalidOiejqResults)
		return
	}
	results := strings.Split(lines[0], " ")
	message := lines[1]
	if results[0] == string(RE) || results[0] == "RV" {
		processInfo.Status = RE
		err = errors.New(message)
		return
	} else if results[0] == string(TLE) {
		processInfo.Status = TLE
		return
	} else if results[0] == string(MLE) {
		processInfo.Status = MLE
		return
	} else if results[0] == string(OLE) {
		processInfo.Status = OLE
		return
	} else if results[0] != string(OK) {
		processInfo.Status = INT
		err = errors.New("invalid oiejq status returned")
		return
	}

	timeMiliseconds, err := strconv.ParseFloat(results[2], 64)
	if err != nil {
		processInfo.Status = INT
		return
	}
	processInfo.TimeInSeconds = timeMiliseconds / 1000
	memory, err := strconv.Atoi(results[4])
	if err != nil {
		processInfo.Status = INT
		return
	}
	processInfo.MemoryInMegabytes = float64(memory) / 1024.0

	processInfo.Status = OK
	return
}

func RunProcessWithOiejq(command string, input io.Reader, oiejqOptions *OiejqOptions) (oiejqProcessInfo ProcessInfo, err error) {
	oiejqResults, err := os.CreateTemp(os.TempDir(), "sio2jail-")
	if err != nil {
		oiejqProcessInfo.Status = INT
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
	processInfo, processErr := RunProcess(oiejqCommand, input, oiejqResults)
	oiejqProcessInfo, err = readOiejqOutput(oiejqResults.Name())
	if oiejqProcessInfo.Status != OK {
		if err != nil && err.Error() == ErrorInvalidOiejqResults {
			oiejqProcessInfo.Status = processInfo.Status
			err = processErr
			return
		}
		return
	}
	if processErr != nil {
		err = processErr
		return
	}
	oiejqProcessInfo.Output = processInfo.Output
	return
}
