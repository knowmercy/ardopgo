git_hash := $(git rev-parse --short HEAD || echo 'development')

# Get current date
current_time = $(date +"%Y-%m-%d:T%H:%M:%S")
#
# # Add linker flags
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_hash}'
#

arch := $(uname -m)

.PHONY:
build:
	@echo "Building binaries..."
	go build -ldflags=${linker_flags} -o=./build/ardop-${arch}-${git_hash} main.go packet.go
	GOOS=linux GOARCH=arm64 go build -ldflags=${linker_flags} -o=./build/ardop-linux-arm64-${git_hash} ./main.go
	GOOS=linux GOARCH=arm go build -ldflags=${linker_flags} -o=./build/ardop-linux-arm-${git_hash} ./main.go
	GOOS=windows GOARCH=amd64 go build -ldflags=${linker_flags} -o=./build/ardop-windows-${git_hash} ./main.go

clean:
	rm -rf ./build/
