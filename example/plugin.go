// This is a sample Drone plugin, patterned loosely after drone-helm.
package main

import (
	"github.com/JaredReisinger/drone-plugin-helper/simple"
)

// GlobalParams are the options available for any/all helm commands
type GlobalParams struct {
	simple.Subcommand
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

// ChartParams common chart-related params
type ChartParams struct {
	CaFile   string
	CertFile string
	KeyFile  string
	Keyring  string
	Password string
	Repo     string
	Username string
	Verify   bool
	Version  string
	Chart    string `cmd:",positional"`
}

// TLSParams are common TLS-related paramters used by severl of the commands.
type TLSParams struct {
	TLS         bool
	TLSCaCert   string
	TLSCert     string
	TLSHostname string
	TLSKey      string
	TLSVerify   bool
}

// CreateParams are the options for "helm create"
type CreateParams struct {
	GlobalParams
	Name    string `cmd:",positional"`
	Starter string
}

// DeleteParams are the options for "helm delete"
type DeleteParams struct {
	GlobalParams
	Description string
	DryRun      bool
	NoHooks     bool
	Purge       bool
	Timeout     int
	TLSParams
	ReleaseName string `cmd:",positional"`
}

// FetchParams are the options for "helm fetch"
type FetchParams struct {
	GlobalParams
	Destination string
	Devel       bool
	Prov        bool
	Untar       bool
	Untardir    string
	ChartParams
}

// GetParams are the options for "helm get"
type GetParams struct {
	GlobalParams
	Revision int32
	TLSParams
	ReleaseName string `cmd:",positional"`
}

// HistoryParams are the options for "helm history"
type HistoryParams struct {
	GlobalParams
	ColWidth uint
	Max      int32
	Output   string
	TLSParams
	ReleaseName string `cmd:",positional"`
}

// InitParams are the options for "helm init"
type InitParams struct {
	GlobalParams
	CanaryImage     bool
	ClientOnly      bool
	DryRun          bool
	ForceUpgrade    bool
	HistoryMax      int
	LocalRepoURL    string
	NetHost         bool
	NodeSelectors   string
	Output          string
	Override        []string
	Replicas        int
	ServiceAccount  string
	SkipRefresh     bool
	StableRepoURL   string
	TillerImage     string
	TillerTLS       bool
	TillerTLSCert   string
	TillerTLSKey    string
	TillerTLSVerify bool
	TLSCaCert       string
	Upgrade         bool
	Wait            bool
}

// InspectParams are the options for "helm inspect"
type InspectParams struct {
	GlobalParams
	ChartParams
}

// InstallParams are the options for "helm install"
type InstallParams struct {
	GlobalParams
	TLSParams
	DepUp        bool
	Description  string
	Devel        bool
	DryRun       bool
	Name         string
	NameTemplate string
	Namespace    string
	NoCrdHook    bool
	NoHooks      bool
	Replace      bool
	Set          []string
	SetFile      []string
	SetString    []string
	Timeout      int
	Values       []string
	Wait         bool
	ChartParams
}

// LintParams are the options for "helm lint"
type LintParams struct {
	GlobalParams
	Namespace string
	Set       []string
	SetFile   []string
	SetString []string
	Strict    bool
	Values    []string
	Path      string `cmd:",positional"`
}

// ListParams are the options for "helm list"
type ListParams struct {
	GlobalParams
	All       bool
	ColWidth  uint
	Date      bool
	Deleted   bool
	Deployed  bool
	Failed    bool
	Max       int
	Namespace string
	Offset    string
	Output    string
	Pending   bool
	Reverse   bool
	Short     bool
	TLSParams
}

func main() {
	simple.ExecSubcommand("helm", map[string]interface{}{
		// "help":    &GlobalParams{},
		"create":  &CreateParams{},
		"delete":  &DeleteParams{},
		"fetch":   &FetchParams{},
		"get":     &GetParams{},
		"history": &HistoryParams{},
		"home":    &GlobalParams{},
		"init":    &InitParams{},
		"inspect": &InspectParams{},
		"install": &InstallParams{},
		"lint":    &LintParams{},
		"list":    &ListParams{},
	})

	// // log.Println("extracting values...")
	// vars := env.Extract(os.Environ(), envPrefix)
	// // log.Printf("extracted: %+v", vars)
	//
	// // log.Println("parsing values...")
	// var params interface{}
	// params = &GlobalParams{}
	// // unused, err := env.Parse(vars, params)
	// _, err := env.Parse(vars, params)
	// if err != nil {
	// 	log.Printf("error: %+v", err)
	// 	return
	// }
	// // log.Printf("parsed: %+v", params)
	// // log.Printf("unused: %+v", unused)
	//
	// // Figure out which actual config to use...
	// // TODO...
	// switch (params.(*GlobalParams)).Command {
	// case "list":
	// 	params = &ListParams{}
	// }
	// //
	// // // re-parse with the specific config...
	// // unused, err := env.Parse(vars, params)
	// // if err != nil {
	// // 	log.Printf("error: %+v", err)
	// // 	return
	// // }
	// // log.Printf("unused: %+v", unused)
	// //
	// // log.Printf("parsed: %+v", params)
	// //
	// // cmd.Exec("helm", params)
	//
	// simple.Exec("helm", params)
}
