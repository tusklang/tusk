ifeq ($(OS),Windows_NT)
	OMM_CMD = ./omm.exe
else
	OMM_CMD = ./omm
endif

all:
	set GOPATH=$(CURDIR)/../
	go build omm.go
	$(OMM_CMD) ./ stdlib/lib.omm -c