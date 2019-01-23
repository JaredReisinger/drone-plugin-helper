package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	// "github.com/JaredReisinger/drone-plugin-helper/env"
)

// Exec runs the command with the given params, exiting the process with any
// errors from the underlying command.  If the command is successful, the
// process is *not* exited; control simply returns from Exec().
func Exec(command string, params interface{}) {
	options, err := Create(params)
	if err != nil {
		log.Fatalf("error creating options: %+v", err)
	}

	err = Run(command, options)
	if err != nil {
		log.Printf("command returned an error: %+v", err)
		exit, ok := err.(*exec.ExitError)
		if ok {
			status, ok := exit.Sys().(syscall.WaitStatus)
			if ok {
				os.Exit(status.ExitStatus())
			}
		}
		log.Fatalln("unable to determine failing exit status")
	}
}

// Run mirrors the os/exec `Cmd.Run()` funciton.
func Run(command string, options []string) (err error) {
	cmd := exec.Command(command, options...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// cmd.Run() doesn't need to quote the arguments because they are already
	// separated into array elements... we would like to show the command-to-run
	// to the user, but don't want to spend the effort to "minimally quote" the
	// way that a hand-typed command-line would be.  Go's "%q" formatting of the
	// array is good enough.
	log.Printf("Running command: %q\n", cmd.Args)
	err = cmd.Run()
	return
}
