VERSION := $(shell (git rev-parse --short HEAD 2>/dev/null|| echo "noversion"))
DIRNAME := $(shell basename $(shell pwd))

version:
	echo ${DIRNAME}
install:
	go install -ldflags '-X github.com/hacker65536/${DIRNAME}/cmd.GitCommit=$(VERSION)'

