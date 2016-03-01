FROM golang:1.6

ENV CGO_ENABLED 0
ENV BUILD_VERSION 0.3

RUN go get -u github.com/golang-devops/yaml-script-runner/...
WORKDIR /go/src/github.com/golang-devops/yaml-script-runner

RUN GOOS=darwin GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-darwin-386
RUN GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-darwin-amd64
RUN GOOS=freebsd GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-freebsd-386
RUN GOOS=freebsd GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-freebsd-amd64
RUN GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-linux-386
RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-linux-amd64
RUN GOOS=netbsd GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-netbsd-386
RUN GOOS=netbsd GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-netbsd-amd64
RUN GOOS=openbsd GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-openbsd-386
RUN GOOS=openbsd GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-openbsd-amd64
RUN GOOS=windows GOARCH=386 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-windows-386.exe
RUN GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -ldflags "-w -s -X main.BuildSha1Version=$BUILD_VERSION" -o=/yaml-script-runner-built/yaml-script-runner-windows-amd64.exe