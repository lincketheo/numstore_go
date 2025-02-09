# Compiler & Flags
CC := gcc
CFLAGS := -Wall -Werror -g -O2
AR := ar rcs

# Source Files
GO_SRC := $(shell find . -name "*.go")
NDBC_LIB := c/libndbc.a

# Binaries
CMD_DIRS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d)
BINARIES := $(notdir $(CMD_DIRS))

all: $(BINARIES)

# Build Go binaries
$(BINARIES): %: $(GO_SRC) $(NDBC_LIB)
	go build -gcflags "all=-N -l" -o $@ ./cmd/$@

# Build static library
$(NDBC_LIB):
	$(MAKE) -C c

# Linting rule
lint:
	golangci-lint run ./...

# Format rule
format:
	find . -name "*.c" -o -name "*.h" | xargs clang-format -i
	gofmt -w $(GO_SRC)

# Clean rule
clean:
	rm -f $(BINARIES) $(NDBC_LIB) $(NDBC_OBJ)
	$(MAKE) -C c clean
