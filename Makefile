GO_SRC		:= $(shell find . -name "*.go")
ndbc_SRC	:= $(wildcard ./c/*.c)
ndbc_OBJ	:= $(ndbc_SRC:.c=.o)
ndbc_LIB	:= c/libndbc.a
CMD_DIRS	:= $(shell find cmd -mindepth 1 -maxdepth 1 -type d)
BINARIES	:= $(notdir $(CMD_DIRS))  

all: $(BINARIES)

$(BINARIES): %: $(GO_SRC) $(ndbc_LIB)
	go build -gcflags "all=-N -l" -o $@ ./cmd/$@

$(ndbc_LIB): $(ndbc_OBJ)
	ar rcs $@ $^

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -f $(BINARIES) $(ndbc_LIB) $(ndbc_OBJ)
