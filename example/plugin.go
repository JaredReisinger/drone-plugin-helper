// This is a sample Drone plugin, patterned loosely after drone-helm.
package main

import (
	"log"
	"os"

	// "github.com/JaredReisinger/drone-plugin-helper/cmd"
	"github.com/JaredReisinger/drone-plugin-helper/env"
	"github.com/JaredReisinger/drone-plugin-helper/simple"
)

const (
	envPrefix string = "PLUGIN_"
)

// GlobalParams are the options available for any/all helm commands
type GlobalParams struct {
	Command                 string `cmd:",positional"`
	Debug                   bool
	Home                    string
	Host                    string
	KubeContext             string
	Kubeconfig              string
	TillerConnectionTimeout int
	TillerNamespace         string

	// lifted from individual commands
	Help bool
}

// ListParams are the options for "helm list"
type ListParams struct {
	GlobalParams
	All         bool
	ColWidth    uint
	Date        bool
	Deleted     bool
	Deployed    bool
	Failed      bool
	Max         int
	Namespace   string
	Offset      string
	Output      string
	Pending     bool
	Reverse     bool
	Short       bool
	TLS         bool
	TLSCaCert   string // golint wants TLSCaCert...
	TLSCert     string
	TLSHostname string
	TLSKey      string
	TLSVerify   bool
}

func main() {
	// log.Println("extracting values...")
	vars := env.Extract(os.Environ(), envPrefix)
	// log.Printf("extracted: %+v", vars)

	// log.Println("parsing values...")
	var params interface{}
	params = &GlobalParams{}
	// unused, err := env.Parse(vars, params)
	_, err := env.Parse(vars, params)
	if err != nil {
		log.Printf("error: %+v", err)
		return
	}
	// log.Printf("parsed: %+v", params)
	// log.Printf("unused: %+v", unused)

	// Figure out which actual config to use...
	// TODO...
	switch (params.(*GlobalParams)).Command {
	case "list":
		params = &ListParams{}
	}
	//
	// // re-parse with the specific config...
	// unused, err := env.Parse(vars, params)
	// if err != nil {
	// 	log.Printf("error: %+v", err)
	// 	return
	// }
	// log.Printf("unused: %+v", unused)
	//
	// log.Printf("parsed: %+v", params)
	//
	// cmd.Exec("helm", params)

	simple.Exec("helm", params)
}
