define newline #define a newline


endef

ifeq ($(OS),Windows_NT)
	BINARY = omm.exe
	COMPILE_ALL = set GOPATH=$(GOPATH_DIR) $(newline) go build omm.go
	CLEAN_CMD = del
	LIB_OAT = stdlib\lib.oat
else
	BINARY = omm
	COMPILE_ALL = @GOPATH=$(GOPATH_DIR) go build omm.go
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
	$(COMPILE_ALL)

$(LIB_OAT):
	-./$(BINARY) ./stdlib ./lib.omm -c
