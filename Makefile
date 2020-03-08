GOPATH := /home/huawei/go:$(PWD)
GOBINARYNAME := network-sniffer
export GOPATH

all: network-sniffer

## lint: Check code style with golint
lint:
	golint

## network-sniffer: Build binary with $(GOBINARYNAME)
network-sniffer: lint
	go build -o $(GOBINARYNAME)

## clean: Clean up
clean:
	rm -rf $(GOBINARYNAME)

## help: Obtain help related information
help: Makefile
	@sed -n 's/^##//p' $<

.PHONY: all network-sniffer
