ifeq ($(OS),Windows_NT)
	BINARY = omml.exe
	LIB_OAT = stdlib\lib.oat
	CLEAN_CMD = del
else
	BINARY = omml
	LIB_OAT = stdlib/lib.oat
	CLEAN_CMD = rm -f
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

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(LIB_OAT)
	-$(CLEAN_CMD) $(BINARY)

.PHONY: $(BINARY)
$(BINARY):
	go build omml.go

.PHONY: lib.oat
lib.oat:
	-./$(BINARY) ./stdlib lib.omm -c
