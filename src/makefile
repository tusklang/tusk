define newline #define a newline


endef

ifeq ($(OS),Windows_NT)
	BINARY = omml.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR) $(newline)
	CLEAN_CMD = del
else
	BINARY = omml
	SET_GOPATH = @GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = rm -f
endif

GOPATH_DIR = $(CURDIR)/../

.PHONY: all
all: $(BINARY) lib.oat

.PHONY: test
test: $(BINARY) lib.oat
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
	$(SET_GOPATH) go build omml.go

.PHONY: lib.oat
lib.oat:
	./$(BINARY) ./stdlib lib.omm -c

.PHONY: gccgo
gccgo:
	$(SET_GOPATH)
	cd lang/types && gccgo -c -o types.o *
	cd oat/helper && gccgo -c -o helper.o *
	cd oat && gccgo -c -o oat.o *
	cd lang/interpreter && gccgo -c -o interpreter.o *
	cd lang/compiler && gccgo -c -o compiler.o *
	gccgo -c omml.go
