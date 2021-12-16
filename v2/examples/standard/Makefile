GIT_COMMIT = $(shell git rev-list -1 HEAD)
BUILD_TIME = $(shell date '+%Y-%m-%d_%H:%M:%S_%Z')

build:
	CGO_ENABLED=0 go build -o standard \
		-ldflags '-X main.gitCommit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)'

run: build
	./standard
