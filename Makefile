.PHONY: clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=fpay
LINUX_386_NAME=$(BINARY_NAME)-linux-386
LINUX_AMD64_NAME=$(BINARY_NAME)-linux-amd64
LINUX_ARM5_NAME=$(BINARY_NAME)-linux-arm5
LINUX_ARM6_NAME=$(BINARY_NAME)-linux-arm6
LINUX_ARM7_NAME=$(BINARY_NAME)-linux-arm7
LINUX_ARM64_NAME=$(BINARY_NAME)-linux-arm64
DARWIN_386_NAME=$(BINARY_NAME)-darwin-386
DARWIN_AMD64_NAME=$(BINARY_NAME)-darwin-amd64
WINDOWS_386_NAME=$(BINARY_NAME)-windows-386
WINDOWS_AMD64_NAME=$(BINARY_NAME)-windows-amd64

all: build
build: fpay
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf pkg
	rm -rf bin
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:

fpay: zlog core account payment confirmation block node

zlog:
	$(GOBUILD) zlog
	$(GOINSTALL) zlog

core:
	$(GOBUILD) fpay/core
	$(GOINSTALL) fpay/core

account:
	$(GOBUILD) fpay/account
	$(GOINSTALL) fpay/account

payment: 
	$(GOBUILD) fpay/payment
	$(GOINSTALL) fpay/payment

confirmation: payment
	$(GOBUILD) fpay/confirmation
	$(GOINSTALL) fpay/confirmation

block: core payment
	$(GOBUILD) fpay/block
	$(GOINSTALL) fpay/block

node:
	$(GOBUILD) fpay/node
	$(GOINSTALL) fpay/node