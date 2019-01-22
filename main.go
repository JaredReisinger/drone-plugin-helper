package main

import (
	"log"
	"os"

	"github.com/JaredReisinger/drone-plugin-helper/cmd"
	"github.com/JaredReisinger/drone-plugin-helper/env"
)

const (
	envPrefix string = "PLUGIN_"
)

func main() {
	log.Println("extracting values...")
	vars := env.Extract(os.Environ(), envPrefix)
	log.Printf("extracted: %+v", vars)

	log.Println("")
	log.Println("parsing values...")
	// cfg := &config{Embedded: &Embedded{}}
	cfg := &config{}
	unused, err := env.Parse(vars, cfg)
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
	log.Printf("parsed: %+v", cfg)
	log.Printf("unused: %+v", unused)

	log.Printf("Inner value: %q (%q)", cfg.Inner, cfg.Embedded.Inner)
	if cfg.IntPtr != nil {
		log.Printf("***IntPtr value: %d", ***cfg.IntPtr)
	} else {
		log.Printf("IntPtr value: <nil>")
	}

	log.Println("")
	log.Println("creating command line...")
	line, err := cmd.Create(cfg)
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
	log.Printf("cmdline: %q", line)
}

type Embedded struct {
	Inner string
}

type config struct {
	Bool bool `cmd:"--bool,no"`

	*Embedded

	IntPtr ***int

	// Int    int    `cmd:"--int"`
	Int8  int8  `cmd:"--int8"`
	Int16 int16 `cmd:"--int16"`
	// Int32  int32  `cmd:"--int32"`
	// Int64  int64  `cmd:"--int64"`
	// Uint   uint   `cmd:"--uint"`
	// Uint8  uint8  `cmd:"--uint8"`
	// Uint16 uint16 `cmd:"--uint16"`
	// Uint32 uint32 `cmd:"--uint32"`
	// Uint64 uint64 `cmd:"--uint64"`

	// String string `cmd:"--string"`

	Default     string
	MissingName string `cmd:""`

	// Map map[string]string
}
