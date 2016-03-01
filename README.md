# yaml-script-runner
Just a basic script runner to easily abort/continue on previous step failure

## Quick Start

### /path/to/my-setup.yaml:

```
variables:
  COMMAND_VAR: which python & echo Hallo
  Machines: ["machine1", "machine2"]

Phase 1: 
  continue_on_failure: true
  inherit_environment: true
  run_parallel: true
  additional_environment:
    - MyEnviron1=Value1
  # executor: ["bash", "-C"] # This
  steps:
    - which python
    - $COMMAND_VAR && echo test123
    - 'repeat::Machines echo {{.}}'
```

### Install and run

```
go get -u github.com/golang-devops/yaml-script-runner
yaml-script-runner "/path/to/my-setup.yaml"
```