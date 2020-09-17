ifeq ($(OS),Windows_NT)
	BINARY = tusk_start.exe
	CLEAN_CMD = del
else
	BINARY = tusk_start
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../

.PHONY: all
all: language

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

.PHONY: language
language:
	go build -a tusk_start.go
