ifeq ($(OS),Windows_NT)
	OMM_CMD = ./omm.exe
	BINARY = omm.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = del
else
	OMM_CMD = ./omm
	BINARY = omm
	SET_GOPATH = @GOPATH=$(GOPATH_DIR)
	CLEAN_CMD = rm -f
endif

GOPATH_DIR := $(CURDIR)/../

all: $(BINARY) stdlib\lib.oat

test: $(BINARY) stdlib\lib.oat
	@echo ---------------------------------- 

	$(OMM_CMD) ./ ./test.omm

clean: $(BINARY) stdlib\lib.oat
	$(CLEAN_CMD) $(BINARY)
	$(CLEAN_CMD) stdlib\lib.oat

$(BINARY): 
	$(SET_GOPATH)
	go build omm.go

stdlib\lib.oat:
	$(OMM_CMD) ./stdlib ./lib.omm -c
