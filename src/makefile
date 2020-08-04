ifeq ($(OS),Windows_NT)
	OMM_CMD = ./omm.exe
	BINARY = omm.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = del
	LIB_OAT = stdlib\lib.oat
else
	OMM_CMD = ./omm
	BINARY = omm
	SET_GOPATH = @GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = rm -f
	LIB_OAT = stdlib/lib.oat
endif

GOPATH_DIR := $(CURDIR)/../

all: $(BINARY) lib.oat

test: $(BINARY) lib.oat
	@echo ------------------------------------
	@echo --------- Start Test File ---------- 
	@echo ------------------------------------
	-$(OMM_CMD) ./ ./test.omm
	@echo ------------------------------------
	@make -s clean_no_echo

clean:
	-$(CLEAN_CMD) $(BINARY)
	-$(CLEAN_CMD) $(LIB_OAT)

.SILENT: clean_no_echo
clean_no_echo:
	-$(CLEAN_CMD) $(BINARY)
	-$(CLEAN_CMD) $(LIB_OAT)

$(BINARY): 
	$(SET_GOPATH)
	go build omm.go

lib.oat:
	$(OMM_CMD) ./stdlib ./lib.omm -c
