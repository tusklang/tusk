ifeq ($(OS),Windows_NT)
	BINARY = omm_start.exe
	CLEAN_CMD = del
else
	BINARY = omm_start
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../

.PHONY: all
all: language

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)
	go mod tidy

.PHONY: language
language:
	go build -a omm_start.go
