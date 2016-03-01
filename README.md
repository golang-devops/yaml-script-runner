# yaml-script-runner
Just a basic script runner to easily abort/continue on previous step failure

## Quick Start

### Install and run

For an example yaml file refer to [example.yml](examples/example.yml)

```
go get -u github.com/golang-devops/yaml-script-runner
yaml-script-runner "/path/to/my-setup.yaml"
```

### Contributions

Refer to the [Godeps.json](Godeps/Godeps.json) file for third-party packages used. But a special thanks to these repositories:
- https://github.com/fatih/color - console colors
- https://github.com/golang-devops/parsecommand - parsing the shell arguments from the single string