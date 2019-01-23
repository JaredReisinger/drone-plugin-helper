package simple

import (
	"log"
	"os"

	"github.com/JaredReisinger/drone-plugin-helper/cmd"
	"github.com/JaredReisinger/drone-plugin-helper/env"
)

// Exec is the all-in-one, "just wrap a command-line tool" method.  If
// you don't need to inspect the values and simply need a one-to-one mapping
// from Drone config through PLUGIN_ environment variables, and into the
// command-line, this is by far the easiest way to get there.
func Exec(command string, params interface{}) {
	_, err := env.Parse(env.Extract(os.Environ(), "PLUGIN_"), params)
	if err != nil {
		log.Fatalf("error parsing environment: %+v\n", err)
	}

	cmd.Exec(command, params)
}

// Subcommand is minimal param data needed to choose subcommand-specific
// parameters.  For convenience, it can also be used as the first embedded
// field in any command-specific parameter struct definitions.  It assumes that
// the subcommand is indicated via the PLUGIN_COMMAND environment variable.
type Subcommand struct {
	Command string `cmd:",positional"`
}

// ExecSubcommand is the all-in-one method for tools which have subcommands,
// like `git` or `helm`.
func ExecSubcommand(command string, paramsMap map[string]interface{}) {
	subcommand := &Subcommand{}
	_, err := env.Parse(env.Extract(os.Environ(), "PLUGIN_"), subcommand)
	if err != nil {
		log.Fatalf("error parsing environment: %+v\n", err)
	}
	params, ok := paramsMap[subcommand.Command]
	if !ok {
		log.Fatalf("subcommand %q not recognized\n", subcommand.Command)
	}

	Exec(command, params)
}
