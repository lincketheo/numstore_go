SRC := $(shell find . -name "*.go")
ndbc_SRC := $(wildcard ./c/*.c)
ndbc_OBJ := $(ndbc_SRC:.c=.o)
ndbc_LIB := c/libndbc.a

ndb: $(SRC) $(ndbc_LIB)
	go build -gcflags "all=-N -l" -o ndb ./cmd/ndb

$(ndbc_LIB): $(ndbc_OBJ)
	ar rcs $@ $^

%.o: %.c
	$(CC) -c $< -o $@

clean:
	rm -f ndb $(ndbc_LIB) $(ndbc_OBJ)
