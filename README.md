# yaml-script-runner
Just a basic script runner to easily abort/continue on previous step failure

## Quick Start

### /path/to/my-setup.yaml:

```
variables:
  COMMAND_VAR: where python && echo Hallo

Phase 1: 
  continue_on_failure: true
  inherit_environment: true
  run_parallel: true
  additional_environment:
    - MyEnviron1=Value1
  executor: ["sh", "-c"]
  steps:
    - where python
    - $COMMAND_VAR && echo test123
```

### Install and run

```
go get -u github.com/golang-devops/yaml-script-runner
yaml-script-runner "/path/to/my-setup.yaml"
```