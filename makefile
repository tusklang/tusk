ifeq ($(OS),Windows_NT)
	BINARY = tusk.exe
	CLEAN_CMD = del
else
	BINARY = tusk.bin
	CLEAN_CMD = rm -f
endif


.PHONY: default
default: all

.PHONY: all
all: pkgs build

.PHONY: build
build:
	go build -o $(BINARY) entry.go

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

.PHONY: pkgs
pkgs:
	go get -u github.com/dlclark/regexp2
	go get -u github.com/llir/llvm
	
.PHONY: prod
prod: all

.PHONY: test
test: build
	cd test && ../$(BINARY)