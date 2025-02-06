# Compiler & Flags
CC := gcc
CFLAGS := -Wall -Werror -g -O2
AR := ar rcs

# Source Files
GO_SRC := $(shell find . -name "*.go")
NDBC_SRC := $(wildcard ./c/*.c)
NDBC_OBJ := $(NDBC_SRC:.c=.o)
NDBC_LIB := libndbc.a

# Binaries
CMD_DIRS := $(shell find cmd -mindepth 1 -maxdepth 1 -type d)
BINARIES := $(notdir $(CMD_DIRS))

.PHONY: all clean lint format

all: $(BINARIES)

# Build Go binaries
$(BINARIES): %: $(GO_SRC) $(NDBC_LIB)
	go build -gcflags "all=-N -l" -o $@ ./cmd/$@

# Build static library
$(NDBC_LIB): $(NDBC_OBJ)
	$(AR) $@ $^

# Compile C source files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

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

