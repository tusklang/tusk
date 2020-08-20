ifeq ($(OS),Windows_NT)
	BINARY = omml.exe
	CLEAN_CMD = del
else
	BINARY = omml
	CLEAN_CMD = rm -f
	SLASHSEP = /
endif

GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: $(BINARY) lib.oat

.PHONY: test
test: all
	@echo ------------------------------------
	@echo --------- Start Test File ----------
	@echo ------------------------------------
	-./$(BINARY) ./ ./test.omm
	@echo ------------------------------------
	@make -s clean_no_echo

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

.SILENT: clean_no_echo
.PHONY: clean_no_echo
clean_no_echo:
	@make -s clean

.PHONY: $(BINARY)
$(BINARY):
	go build omml.go

.PHONY: lib.oat
lib.oat:
	./$(BINARY) ./stdlib lib.omm -c
