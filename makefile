ifeq ($(OS),Windows_NT)
	BINARY = omml.exe
	CLEAN_CMD = del
	FFIGENERATE = ffi\generate\gen.exe
	FFICALLSWITCH = ffi\callswitch.h
else
	BINARY = omml
	CLEAN_CMD = rm -f
	SLASHSEP = /
	FFIGENERATE = ffi/generate/gen
	FFICALLSWITCH = ffi/callswitch.h
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
	-$(CLEAN_CMD) $(FFIGENERATE)
	-$(CLEAN_CMD) $(FFICALLSWITCH)
	-$(CLEAN_CMD) $(BINARY)

.SILENT: clean_no_echo
.PHONY: clean_no_echo
clean_no_echo:
	@make -s clean

.PHONY: $(BINARY)
$(BINARY):
	g++ ffi/generate/wingenerate_argswitch.cc -o $(FFIGENERATE)
	./$(FFIGENERATE)
	go build omml.go

.PHONY: lib.oat
lib.oat:
	./$(BINARY) ./stdlib lib.omm -c
