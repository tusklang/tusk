ifeq ($(OS),Windows_NT)
	BINARY = omm_start.exe
	CLEAN_CMD = del
else
	BINARY = omm_start
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
	go build omm_start.go
