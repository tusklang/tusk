ifeq ($(OS),Windows_NT)
	OMM_CMD = ./omm.exe
	BINARY = omm.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = del
	LIB_OAT = stdlib\lib.oat
	MANGO_OAT = mango\mango.oat
else
	OMM_CMD = ./omm
	BINARY = omm
	SET_GOPATH = @GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = rm -f
	LIB_OAT = stdlib/lib.oat
	MANGO_OAT = mango/mango.oat
endif

GOPATH_DIR := $(CURDIR)/../

all: $(BINARY) $(LIB_OAT) $(MANGO_OAT)

test: $(BINARY) $(LIB_OAT)
	@echo ------------------------------------
	@echo --------- Start Test File ---------- 
	@echo ------------------------------------
	-$(OMM_CMD) ./ ./test.omm
	@echo ------------------------------------
	@make -s clean_no_echo

clean:
	-$(CLEAN_CMD) $(BINARY)
	-$(CLEAN_CMD) $(LIB_OAT)
	-$(CLEAN_CMD) $(MANGO_OAT)

.SILENT: clean_no_echo
clean_no_echo:
	@make -s clean

$(BINARY): 
	$(SET_GOPATH)
	go build omm.go

$(LIB_OAT):
	$(OMM_CMD) ./stdlib ./lib.omm -c

$(MANGO_OAT):
	$(OMM_CMD) ./mango ./mango.omm -c
