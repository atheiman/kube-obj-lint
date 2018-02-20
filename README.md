# Kubernetes Object Validation with Go

I'm learning Go and Kubernetes, I created a simple Go tool that could be used to validate logic in Kubernetes object yaml files.

So far I've only created a couple logic checks around Pod volumes:

- if a container tries to mount a volume, the volume must be defined in the volumes spec
  - `error` level - fails validation
- if a volume is defined in the volumes spec but not mounted into any containers
  - `warn` level - validation still passes

## Usage

See the Travis tests to see how to run this.

## To Do

Future things I could do to learn more:

- use [a more sophisticated logger](https://github.com/op/go-logging) for better log output (`warn` vs `error`)
- handle more object types, determine based on yaml content
- accept a path to search for yaml files in (`files/*.yml`, `.`)
- make this an installable go package that builds in travis and can be downloaded from a url and run as a binary directly
