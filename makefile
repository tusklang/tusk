ifeq ($(OS),Windows_NT)
	BINARY = tusk.exe
	CLEAN_CMD = del
else
	SET_GOPATH = GOPATH=$(GOPATH)
	BINARY = tusk.bin
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: default
default: all

.PHONY: all
all: pkgs
	$(SET_GOPATH) go build -o $(BINARY) entry.go

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

pkgs:
	go get github.com/dlclark/regexp2

prod:

test: all
	./$(BINARY) -wd=./test/