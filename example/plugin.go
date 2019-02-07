// This is a sample Drone plugin, patterned loosely after drone-helm.
package main

import (
	"github.com/JaredReisinger/drone-plugin-helper/simple"
)

// GlobalParams are the options available for any/all helm commands
type GlobalParams struct {
	simple.Command
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

// DependencyBuildParams are the options for "helm dependency build"
type DependencyBuildParams struct {
	GlobalParams
	Keyring string
	Verify  bool
	Chart   string `cmd:",positional"`
}

// DependencyListParams are the options for "helm dependency list"
type DependencyListParams struct {
	GlobalParams
	Chart string `cmd:",positional"`
}

// DependencyUpdateParams are the options for "helm dependency update"
type DependencyUpdateParams struct {
	GlobalParams
	Keyring     string
	SkipRefresh bool
	Verify      bool
	Chart       string `cmd:",positional"`
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

// GetParams are the options for "helm get" (and subcommands)
type GetParams struct {
	GlobalParams
	Revision int32
	TLSParams
	ReleaseName string `cmd:",positional"`
}

// GetValuesParams are the options for "helm get values"
type GetValuesParams struct {
	All bool
	GetParams
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

// PackageParams are the options for "helm package"
type PackageParams struct {
	GlobalParams
	AppVersion       string
	DependencyUpdate bool
	Destination      string
	Key              string
	Keyring          string
	Save             bool // true by default
	Sign             bool
	Version          string
}

// PluginInstallParams are the options for "helm plugin install"
type PluginInstallParams struct {
	GlobalParams
	Path    string `cmd:",positional"`
	Version string
}

// PluginParams are the options for "helm plugin remove" and
// "helm plugin update"
type PluginParams struct {
	GlobalParams
	Plugins []string `cmd:",positional"`
}

// RepoAddParams are the options for "helm repo add"
type RepoAddParams struct {
	GlobalParams
	CaFile   string
	CertFile string
	KeyFile  string
	NoUpdate bool
	Password string
	Username string
	Name     string `cmd:",positional"`
	URL      string `cmd:",positional"`
}

// RepoIndexParams are the options for "helm repo index"
type RepoIndexParams struct {
	GlobalParams
	Merge     string
	URL       string
	Directory string `cmd:",positional"`
}

// RepoRemoveParams are the options for "helm repo remove"
type RepoRemoveParams struct {
	GlobalParams
	Name string `cmd:",positional"`
}

// ResetParams are the options for "helm reset"
type ResetParams struct {
	GlobalParams
	Force          bool
	RemoveHelmHome bool
	TLSParams
}

// RollbackParams are the options for "helm rollback"
type RollbackParams struct {
	GlobalParams
	Description  string
	DryRun       bool
	Force        bool
	NoHooks      bool
	RecreatePods bool
	Timeout      int
	TLSParams
	Wait     bool
	Release  string `cmd:",positional"`
	Revision string `cmd:",positional"`
}

// SearchParams are the options for "helm search"
type SearchParams struct {
	GlobalParams
	ColWidth uint
	Regexp   bool
	Version  string
	Versions bool
	Keyword  string `cmd:",positional"` // first?
}

// ServeParams are the options for "helm serve"
type ServeParams struct {
	GlobalParams
	Address  string
	RepoPath string
	URL      string
}

// StatusParams are the options for "helm status"
type StatusParams struct {
	GlobalParams
	Output   string
	Revision int32
	TLSParams
	ReleaseName string `cmd:",positional"`
}

// TemplateParams are the options for "helm template"
type TemplateParams struct {
	GlobalParams
	Execute      []string
	IsUpgrade    bool
	KubeVersion  string
	Name         string
	NameTemplate string
	Notes        bool
	OutputDir    string
	Set          []string
	SetFile      []string
	SetString    []string
	Values       []string
	Chart        string `cmd:",positional"`
}

// TestParams are the options for "helm test"
type TestParams struct {
	GlobalParams
	Cleanup bool
	Timeout int
	TLSParams
	Release string `cmd:",positional"`
}

// UpgradeParams are the options for "helm upgrade"
type UpgradeParams struct {
	GlobalParams
	TLSParams
	Description  string
	Devel        bool
	DryRun       bool
	Force        bool
	Install      bool
	Namespace    string
	NoHooks      bool
	RecreatePods bool
	ResetValues  bool
	ReuseValues  bool
	Set          []string
	SetFile      []string
	SetString    []string
	Timeout      int
	Values       []string
	Wait         bool
	Release      string `cmd:",positional"`
	ChartParams
}

// VerifyParams are the options for "helm verify"
type VerifyParams struct {
	GlobalParams
	Keyring string
	Path    string `cmd:",positional"`
}

// VersionParams are the options for "helm version"
type VersionParams struct {
	GlobalParams
	Client   bool
	Server   bool
	Short    bool
	Template string
	TLSParams
}

func main() {
	simple.ExecCommand("helm", map[string]interface{}{
		// 'help' is a bogus command that 'helm' responds to with basic usage.
		"help": &GlobalParams{},

		"create":            &CreateParams{},
		"delete":            &DeleteParams{},
		"dependency build":  &DependencyBuildParams{},
		"dependency list":   &DependencyListParams{},
		"dependency update": &DependencyUpdateParams{},
		"fetch":             &FetchParams{},
		"get":               &GetParams{},
		"get hooks":         &GetParams{},
		"get manifest":      &GetParams{},
		"get values":        &GetValuesParams{},
		"history":           &HistoryParams{},
		"home":              &GlobalParams{},
		"init":              &InitParams{},
		"inspect":           &InspectParams{},
		"inspect chart":     &InspectParams{},
		"inspect readme":    &InspectParams{},
		"inspect values":    &InspectParams{},
		"install":           &InstallParams{},
		"lint":              &LintParams{},
		"list":              &ListParams{},
		"package":           &PackageParams{},
		"plugin install":    &PluginInstallParams{},
		"plugin list":       &GlobalParams{},
		"plugin remove":     &PluginParams{},
		"plugin update":     &PluginParams{},
		"repo add":          &RepoAddParams{},
		"repo index":        &RepoIndexParams{},
		"repo list":         &GlobalParams{},
		"repo remove":       &RepoRemoveParams{},
		"repo update":       &GlobalParams{},
		"reset":             &ResetParams{},
		"rollback":          &RollbackParams{},
		"search":            &SearchParams{},
		"serve":             &ServeParams{},
		"status":            &StatusParams{},
		"template":          &TemplateParams{},
		"test":              &TestParams{},
		"upgrade":           &UpgradeParams{},
		"verify":            &VerifyParams{},
		"version":           &VersionParams{},
	})
}
