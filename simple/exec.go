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
