# drone-plugin-helper


## Background

It appears that a ridiculous number of Drone plugins follow the same basic pattern if they are encapsulating an existing command-line tool... and a _lot_ of plugins do this.

  * Drone ensures that the plugin settings from the `.drone.yml` file are exposed as environment variables.
  * A generic CLI library (often urfave/cli) is used to map all manner of CLI+environment to a Config struct.
  * A bunch of repetitive code is written to remap the Config into the underlying tool's actual command-line.

Rather than every plugin re-creating the env-to-config-to-commandline processing again and again, there should be a library that does this basic work _once_, with minimal redundancy for plugin authors.  This project attempts to do just that.


## Concepts


### Environment to config mapping

Since a Drone plugin _only_ gets it input/settings via environment variables, a general-purpose CLI helper is overkill, and adds unneeded complexity.  We should be able to deserialize the environment variables directly into a struct, based solely on the names and types.  (_Maybe_ with tagged members for additional options?)

1-to-1 name mapping, to reduce typing and confusion?  Drone upper-cases the variables and prefixes with `PLUGIN_` so it should be easy to find them, and detect extras...


### Config to command-line mapping

Similarly, the vast majority have a 1-to-1 mapping between a config value and the underlying tool's command line.  Simple struct member tagging should allow the command-line to be almost directly serialzed from the struct.
