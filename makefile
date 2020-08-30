ifeq ($(OS),Windows_NT)
	BINARY = omm.exe
	LIB_OAT = stdlib\lib.oat
	CLEAN_CMD = del
	GOATV_SHARED_DYNAMIC_LIB = goatv.dll
else
	BINARY = omm
	LIB_OAT = stdlib/lib.oat
	CLEAN_CMD = rm -f
	GOATV_SHARED_DYNAMIC_LIB = goatv.so
endif

GOATV_HEADER = goatv.h
GOPATH = $(CURDIR)/../../../../

.PHONY: all
all: $(BINARY) lib.oat goatv

.PHONY: test
test: all
	@echo ------------------------------------
	@echo --------- Start Test File ----------
	@echo ------------------------------------
	-./$(BINARY) ./test.omm
	@echo ------------------------------------

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(LIB_OAT)
	-$(CLEAN_CMD) $(BINARY)
	-$(CLEAN_CMD) $(GOATV_SHARED_DYNAMIC_LIB)
	-$(CLEAN_CMD) $(GOATV_HEADER)

.PHONY: $(BINARY)
$(BINARY):
	go build omm.go

.PHONY: lib.oat
lib.oat:
	-./$(BINARY) stdlib/lib.omm -c

.PHONY: goatv
goatv:
	go build -buildmode=c-shared -o $(GOATV_SHARED_DYNAMIC_LIB) oat/goatv/goatv.go
