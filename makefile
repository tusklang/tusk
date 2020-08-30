ifeq ($(OS),Windows_NT)
	BINARY = omm.exe
	CLEAN_CMD = del
else
	BINARY = omm
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: language

.PHONY: test
test: all
	@echo ------------------------------------
	@echo --------- Start Test File ----------
	@echo ------------------------------------
	-./$(BINARY) ./test.omm
	@echo ------------------------------------

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

.PHONY: language
language:
	go build omm.go
