ifeq ($(OS),Windows_NT)
	BINARY = tuskstart.exe
	CLEAN_CMD = del
else
	BINARY = tuskstart
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: language go.mod

.PHONY: clean
clean:
	-go mod tidy
	-$(CLEAN_CMD) $(BINARY)

go.mod:
	go mod init

.PHONY: language
language:
	go get -u
	go build -a tuskstart.go