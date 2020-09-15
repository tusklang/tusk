ifeq ($(OS),Windows_NT)
	BINARY = ka_start.exe
	CLEAN_CMD = del
else
	BINARY = ka_start
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
	go build ka_start.go
