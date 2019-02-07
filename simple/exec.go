package simple

import (
	"fmt"
	"log"
	"os"

	"github.com/JaredReisinger/drone-plugin-helper/cmd"
	"github.com/JaredReisinger/drone-plugin-helper/env"
)

// TODO: better name than "Exec()"?  That doesn't imply any of the parsing that
// will occur.  Perhaps "Handoff()", or "Passthrough()"?

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

// Command is minimal param data needed to choose command-specific parameters.
// For convenience, it can also be used as the first embedded field in any
// command-specific parameter struct definitions.  It assumes that the command
// is indicated via the PLUGIN_COMMAND environment variable, and any subcommand
// via PLUGIN_SUBCOMMAND.
type Command struct {
	Command    string `cmd:",positional"`
	Subcommand string `cmd:",positional"`
}

// ExecCommand is the all-in-one method for tools which have subcommands,
// like `git` or `helm`.
func ExecCommand(command string, paramsMap map[string]interface{}) {
	commandParams := &Command{}
	_, err := env.Parse(env.Extract(os.Environ(), "PLUGIN_"), commandParams)
	if err != nil {
		log.Fatalf("error parsing environment: %+v\n", err)
	}
	key := commandParams.Command
	if commandParams.Subcommand != "" {
		key = fmt.Sprintf("%s %s", commandParams.Command, commandParams.Subcommand)
	}

	params, ok := paramsMap[key]
	if !ok {
		log.Fatalf("command %q not recognized\n", key)
	}

	Exec(command, params)
}
