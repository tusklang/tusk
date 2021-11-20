ifeq ($(OS),Windows_NT)
	BINARY = tusk.exe
	CLEAN_CMD = del
else
	BINARY = tusk.bin
	CLEAN_CMD = rm -f
endif


.PHONY: default
default: all

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY) entry.go

.PHONY: clean
clean:
	-$(CLEAN_CMD) $(BINARY)

.PHONY: prod
prod: all

.PHONY: test
test: build
	cd test && ../$(BINARY)
	cd test && clang test.ll -o test.bin
	@echo ""
	@echo "Running Test"
	@echo "----------------------"
	@echo ""
	cd test && ./test.bin