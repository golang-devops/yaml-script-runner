# yaml-script-runner
Just a basic script runner to easily abort/continue on previous step failure

## Quick Start

### Install and run

#### Install on linux with curl script

Ensure first that the `MY_OS`, `MY_ARCH` and `VER` variables suit your system. The script does not validate the arguments but just print out a probable list.

This list is directly based on the available download files found here: https://github.com/golang-devops/yaml-script-runner/releases/latest.


```
MY_OS=linux
MY_ARCH=amd64
VER=v0.4
curl -s https://raw.githubusercontent.com/golang-devops/yaml-script-runner/master/_scripts/install.sh | sudo bash /dev/stdin $VER $MY_OS $MY_ARCH
```

#### Install from source (golang must be installed)

For an example yaml file refer to [example.yml](examples/example.yml)

```
go get -u github.com/golang-devops/yaml-script-runner
yaml-script-runner "/path/to/my-setup.yaml"
```

### Contributions

Refer to the [Godeps.json](Godeps/Godeps.json) file for third-party packages used. But a special thanks to these repositories:
- https://github.com/fatih/color - console colors
- https://github.com/golang-devops/parsecommand - parsing the shell arguments from the single string