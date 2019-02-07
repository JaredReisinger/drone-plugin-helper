# plugin-example

A case-study for `drone-plugin-helper`.


## How to build

This example directory contains the (surprisingly small) source code for a Drone plugin for Helm, and a `Dockerfile` to build it into a fairly minimal image.  You can try it out yourself with:

```bash
docker build -t plugin-example:0.1 .
docker run --rm --env 'PLUGIN_COMMAND=list' --env 'PLUGIN_HELP=true' plugin-example:0.1
```

You can also build/test locally (not through Docker) with:

```bash
go build -v . -o plugin-example
PLUGIN_COMMAND=list PLUGIN_HELP=true ./plugin-example
```


## Behind the scenes

This Helm wrapper makes liberal use of the [`drone-plugin-helper/simple`](../simple/) package’s ability to handle subcommand-based command-lines.  The vast majority of the code in `plugin.go` (86.7%) is just to define the command-line parameters.  A handful of lines maps each command/sub-command to the appropriate parameter struct, and invoking Helm itself is really just a single line.

Here’s the breakdown:

| part                        | lines | percentage |
|-----------------------------|------:|-----------:|
| package/imports             |     4 |       1.5% |
| param definitions           |   307 |      86.7% |
| command-to-param mapping    |    42 |      11.9% |
| command-line tool execution |     1 |       0.3% |
| _**total**_                 |  _354_|    _100.0%_|

With `drone-plugin-helper` in place, the majority of the effort is in defining the parameters specific to the command-line tool—the things that make _this_ tool different from another tool—and that’s as it should be.
