ifeq ($(OS),Windows_NT)
	BINARY = tuskstart.exe
	CLEAN_CMD = del
else
	BINARY = tuskstart
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: language

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)
.PHONY: language
language:
	go build -a tuskstart.go