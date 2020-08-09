define newline #define a newline


endef

ifeq ($(OS),Windows_NT)
	BINARY = omm.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR) $(newline)
	CLEAN_CMD = del
	LIB_OAT = stdlib\lib.oat
else
	BINARY = omm
	SET_GOPATH = export GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = rm -f
	LIB_OAT = stdlib/lib.oat
endif

GOPATH_DIR = $(CURDIR)/../

.PHONY: all
all: $(BINARY) $(LIB_OAT)

.PHONY: test
test: $(BINARY) $(LIB_OAT)
	@echo ------------------------------------
	@echo --------- Start Test File ---------- 
	@echo ------------------------------------
	-./$(BINARY) ./ ./test.omm
	@echo ------------------------------------
	@make -s clean_no_echo

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)
	-$(CLEAN_CMD) $(LIB_OAT)

.SILENT: clean_no_echo
.PHONY: clean_no_echo
clean_no_echo:
	@make -s clean

$(BINARY):
	$(SET_GOPATH)
	go build omm.go

$(LIB_OAT):
	-./$(BINARY) ./stdlib ./lib.omm -c

.PHONY: gccgo
gccgo:
	$(SET_GOPATH)
	cd lang/types && gccgo -c -o types.o *
	cd oat/helper && gccgo -c -o helper.o *
	cd oat && gccgo -c -o oat.o *
	cd lang/interpreter && gccgo -c -o interpreter.o *
	cd lang/compiler && gccgo -c -o compiler.o *
	gccgo -c omm.go
