# drone-plugin-helper

Simplifies Drone plugins that wrap command-line tools.


## Getting started

To write a Drone plugin to wrap a command-line tool, all you really need to do is create a struct that represents the available command-line options, and then let the library take care of everything else.  For example, to handle the first several options for the `curl` command:

```text
Usage: curl [options...] <url>
     --abstract-unix-socket <path> Connect via abstract Unix domain socket
     --anyauth       Pick any authentication method
 -a, --append        Append to target file when uploading
     --basic         Use HTTP Basic Authentication
     --cacert <file> CA certificate to verify peer against
     --capath <dir>  CA directory to verify peer against
 -E, --cert <certificate[:password]> Client certificate file and password
     --cert-status   Verify the status of the server certificate
     --cert-type <type> Certificate file type (DER/PEM/ENG)
```

The code is as simple as:

```Go
package main

import (
  "github.com/JaredReisinger/drone-plugin-helper/simple"
)

type Params struct {
  AbstractUnixSocket string
  Anyauth            bool
  Append             bool
  Basic              bool
  Cacert             string
  Capath             string
  Cert               string
  CertStatus         bool
  CertType           string
}

func main() {
  simple.Exec("curl", &Params{})
}
```

All that’s left is for you to package the tool along with the Go executable into a Docker image.

Because of the way that Drone maps settings to environment variables, and the way that `drone-plugin-helper` maps these names to Go member fields, your `drone.yml` file could then specify (for example):

```yaml
some_curl_step:
  image: drone-curl
  anyauth: true
  cacert: ./some/cert
  cert_status: true
  cert_type: PEM
```

... and `drone-plugin-helper` would generate the final command-line:

```sh
curl --anyauth --cacert ./some/cert --cert-status --cert-type PEM
```

Behind the scenes, Drone maps the settings to the environment variables:

```text
PLUGIN_ANYAUTH=true
PLUGIN_CACERT=./some/cert
PLUGIN_CERT_STATUS=true
PLUGIN_CERT_TYPE=PEM
```

... and `drone_plugin_helper` matches these with the `Anyauth`, `Cacert`, `CertStatus`, and `CertType` members.


### “Subcommand”-style tools

Another very common pattern for command-line tools is for the initial command-line argument to be a subcommand, often with its own specific options.  Tools like `git` and `helm` are examples of this.  Since this is such a common pattern, there is a helper for this, as well.  The [`example/`](./example/) subdirectory uses this helper to show how a plugin for `helm` could be written. (See [Example / Case study](#example--case-study) below for further information.)


### Overriding the defaults

If the helpers’ default processing doesn’t meet your needs, you can tag the fields with `env:""` and/or `cmd:""` metadata:

```Go
type Params struct {
  WeirdName string `cmd:"--surprise"`   // use "--surprise" instead of "--weird-name" as the option name
  Command   string `cmd:",positional"`  // use the value directly, with no "--command" flag prefix
  Extra     string `cmd:",omit"`
}
```

But please see “Best practices”, below, for ways to avoid needing these overrides.


### More-complex handling

While the behavior of [`drone-plugin-helper/simple`](./simple/) should handle the vast majority of cases, feel free to use the [`/env`](./env/) or [`/cmd`](./cmd/) packages directly if you need to add your own logic in between the environment variable parsing and the command-line generation.  You may find that [`/env`](./env/) alone is a simpler way to expose your plugin’s parameters even if you’re not wrapping an underlying command-line tool.


## Best practices

### Work “bottom-up”

The ultimate goal is to generate a command-line for the underlying tool, so it makes sense to choose struct member names and types that facilitate this.  Doing so will reduce the need for struct field tag metadata to “fix” the command-line option names.  The `drone-plugin-helper/cmd` methods were designed to generate the “expected” command-line option name based on the Go name: the field `Basic` generates an option named `--basic`, and the field `CertStatus` generates `--cert-status`.

The rule of thumb in naming a Go member for a command-line parameter is to capitalize the first letter of each hyphen-separated term, and then remove the hyphens: `--cert-status` ⇒ `--Cert-Status` ⇒ `CertStatus`.  The helpers are aware of Go's linting rules about capitalizing certain acronyms and respects them.  For example, the proper Go member name for `--tls-cert` is `TLSCert` (not `TlsCert`).  Similar logic in `drone-plugin-helper/env` will look for environment variables with the equivalent environment name: a `TLSCert` member looks for the `PLUGIN_TLS_CERT` environment variable.


### Use pointers if “zero values” are valid options

In typical usage, you will rarely need to inspect the values in the Go struct at all; they become simple pass-throughs from the Drone environment variables to the underlying tool's command-line.  The `drone-plugin-helper` tools will automatically handle struct fields that are pointers: automatically creating and dereferencing as needed.  If a field is *not* a pointer, the option is only emitted when it's not a "zero value" for the type (that is, int options of `0` are not emitted, nor are empty strings).  If a "zero value" has meaning for the underlying tool, you can use a pointer to the type instead; it will only be allocated if an environment value is provided, and any non-`nil` value will be emitted as a command-line option.

For example, if the Go struct is defined as:

```Go
type Params struct {
  Example1 int
  Example2 *int
}
```

then if the environment variables are `PLUGIN_EXAMPLE1=0` and `PLUGIN_EXAMPLE2=0`, only the command-line option `--example2 0` would be created.  The `--example1` option is not created because the value is the "zero value" for the field.

> NOTE: Would it be better to recommend "use pointers by default" because it more accurately passes along the intent of the caller?


## Example / Case study

As a case study, see [`example/`](./example/) to see how `drone-plugin-helper`could be used to create a plugin for Helm.


----

## Background

There are a huge number of Drone plugins that are simply wrappers around an existing command-line tool.  This should be no surprise, as the typical Un*x toolchain is a set of several command-line tools!  The plugins all follow the same overall pattern:

  * Drone ensures that the plugin settings from the `.drone.yml` file are exposed as environment variables.

  * A generic CLI library (often urfave/cli) is used to map all manner of CLI+environment to a struct whose members support the underlying tool’s options.

  * A bunch of repetitive code is written to remap the struct into the underlying tool's actual command-line.

Rather than every plugin re-creating the env-to-struct-to-commandline processing again and again, there should be a library that does this basic work _once_, with minimal redundancy for plugin authors.  This project attempts to do just that.


## Notes

### Environment to config mapping

Since a Drone plugin _only_ gets it input/settings via environment variables, a general-purpose CLI helper is overkill, and adds unneeded complexity.  We should be able to deserialize the environment variables directly into a struct, based solely on the names and types.  (_Maybe_ with tagged members for additional options?)

1-to-1 name mapping, to reduce typing and confusion?  Drone upper-cases the variables and prefixes with `PLUGIN_` so it should be easy to find them, and detect extras...


### Config to command-line mapping

Similarly, the vast majority have a 1-to-1 mapping between a config value and the underlying tool's command line.  Simple struct member tagging should allow the command-line to be almost directly serialzed from the struct.



## TODO (ideas)

* [X] ~~handle "well-known" all-caps abbreviations "TLS", "XML", etc., so that they are parsed correctly: `TLSXMLInfo` -> `TLS` `XML` `Info` > `--tls-xml-info`.  Similarly, `PLUGIN_TLS_XML_INFO` should become `TLSXMLInfo`, not `TlsXmlInfo`.~~

* [X] ~~put the current `cmd.DronePlugin()` into a separate package for "all-in-one" tools.  Perhaps `allInOne` or `simple`?~~ `simple.Exec()`

* [X] ~~add (to `simple` as previous?) a helper for typical "subcommand" behavior, a la `helm`.  This would allow an easy "map the initial command to a param set" which is another 80/20 case.~~

* [ ] helper that can generate a stub doc showing the supported environment vsriables and the allowed syntax... and the mapping to eventual command-line?
