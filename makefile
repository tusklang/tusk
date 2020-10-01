ifeq ($(OS),Windows_NT)
	BINARY = tuskstart.exe
	SYSC_OBJECT = native\systables\sys.syso
	CLEAN_CMD = del
else
	BINARY = tuskstart.out
	SYSC_OBJECT = native/systables/sys.syso
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: default
default: all

.PHONY: all
all: language

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(SYSC_OBJECT)
	-$(CLEAN_CMD) $(BINARY)

.PHONY: language
language:
	gcc --std=c99 -o $(SYSC_OBJECT) -c native/systables/syscalls/sys.c
	go build -a -o $(BINARY) tuskstart.go