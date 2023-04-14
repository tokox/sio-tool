package cmd

import (
	"os"
	"sio-tool/codeforces_client"
	"sio-tool/config"
	"sio-tool/szkopul_client"
	"strings"

	"github.com/docopt/docopt-go"
)

// ParsedArgs parsed arguments
type ParsedArgs struct {
	CodeforcesInfo codeforces_client.Info
	SzkopulInfo    szkopul_client.Info
	File           string
	Generator      string
	Solve          string
	Brute          string
	Specifier      []string `docopt:"<specifier>"`
	Alias          string   `docopt:"<alias>"`
	Accepted       bool     `docopt:"ac"`
	All            bool     `docopt:"all"`
	Handle         string   `docopt:"<handle>"`
	Version        string   `docopt:"{version}"`
	Config         bool     `docopt:"config"`
	Submit         bool     `docopt:"submit"`
	List           bool     `docopt:"list"`
	Parse          bool     `docopt:"parse"`
	Gen            bool     `docopt:"gen"`
	Test           bool     `docopt:"test"`
	Watch          bool     `docopt:"watch"`
	Open           bool     `docopt:"open"`
	Stand          bool     `docopt:"stand"`
	Sid            bool     `docopt:"sid"`
	Race           bool     `docopt:"race"`
	Pull           bool     `docopt:"pull"`
	Clone          bool     `docopt:"clone"`
	Upgrade        bool     `docopt:"upgrade"`
	StressTest     bool     `docopt:"stress-test"`
	Codeforces     bool
	Szkopul        bool
}

// Args global variable
var Args *ParsedArgs

func SubPath(parent, sub string) bool {
	return strings.HasPrefix(parent, sub)
}

func determineClient() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg := config.Instance
	codeforcesDir := SubPath(path, cfg.FolderName["codeforces-root"])
	if codeforcesDir {
		Args.Codeforces = true
		return nil
	}
	szkopulDir := SubPath(path, cfg.FolderName["szkopul-root"])
	if szkopulDir {
		Args.Szkopul = true
		return nil
	}
	return nil
}

func parseArgs(opts docopt.Opts) error {
	if file, ok := opts["--file"].(string); ok {
		Args.File = file
	} else if file, ok := opts["<file>"].(string); ok {
		Args.File = file
	}
	determineClient()
	if Args.Codeforces {
		return parseArgsCodeforces(opts)
	}
	if Args.Szkopul {
		return parseArgsSzkopul(opts)
	}
	return nil
}
