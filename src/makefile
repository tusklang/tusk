ifeq ($(OS),Windows_NT)
	OMM_CMD = ./omm.exe
	SET_GOPATH = set GOPATH=$(GOPATH_DIR)
else
	OMM_CMD = ./omm
	SET_GOPATH = @GOPATH=$(GOPATH_DIR)
endif

GOPATH_DIR := $(CURDIR)/../

build: omm.go
	$(SET_GOPATH)
	go build omm.go
	$(OMM_CMD) ./stdlib ./lib.omm -c

test: omm.go
	$(SET_GOPATH)
	go build omm.go
	$(OMM_CMD) ./stdlib ./lib.omm -c

	@echo ---------------------------------- 

	$(OMM_CMD) ./ ./test.omm
