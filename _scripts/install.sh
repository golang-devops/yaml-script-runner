#!/bin/sh

VERSION=v0.3
allowed_os=( "darwin" "freebsd" "linux" "netbsd" "openbsd" )
allowed_arch=( "386" "amd64" )
 
os=$1
arch=$2
 
## Do not validate for now, otherwise we might prevent a download that actually exists if this script is old
# valid_os=0
# for i in "${allowed_os[@]}"
# do
#     if [ "$os" == "$i" ]; then
#         valid_os=1
#     fi
# done
# 
# valid_arch=0
# for i in "${allowed_arch[@]}"
# do
#     if [ "$arch" == "$i" ]; then
#         valid_arch=1
#     fi
# done
# 
# if [ $valid_os != 1 ]; then
#     echo Invalid os "$os" - valid list is ${allowed_os[@]}
#     exit 1
# fi
# 
# if [ $valid_arch != 1 ]; then
#     echo Invalid arch "$arch" - valid list is ${allowed_arch[@]}
#     exit 1
# fi



echo Probable list of OS is ${allowed_os[@]}
echo Probable list of ARCH is ${allowed_arch[@]}
echo For a full list of supported OS-ARCH look at list in https://github.com/golang-devops/yaml-script-runner/releases/latest

echo Your input OS = "$os" and ARCH = "$arch"
echo This is an install script for yaml-script-runner version $VERSION

URL=https://github.com/golang-devops/yaml-script-runner/releases/download/$VERSION/$os_$arch_yaml-script-runner

echo Now fetching binary at url $URL
wget $URL -O /usr/local/bin/yaml-script-runner
chmod +x /usr/local/bin/yaml-script-runner