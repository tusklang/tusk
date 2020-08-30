ifeq ($(OS),Windows_NT)
	BINARY = omm.exe
	LIB_OAT = stdlib\lib.oat
	CLEAN_CMD = del
else
	BINARY = omm
	LIB_OAT = stdlib/lib.oat
	CLEAN_CMD = rm -f
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: language stdlib goatv

.PHONY: test
test: all stdlib
	@echo ------------------------------------
	@echo --------- Start Test File ----------
	@echo ------------------------------------
	-./$(BINARY) ./test.omm
	@echo ------------------------------------

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(LIB_OAT)
	-$(CLEAN_CMD) $(BINARY)

.PHONY: language
language:
	go build omm.go

.PHONY: stdlib
stdlib:
	-./$(BINARY) stdlib/lib.omm -c
