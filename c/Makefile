INCLUDE_DIR := ./include
SRC_DIR := ./src
APPS_DIR := ./apps

CC := gcc
CFLAGS := -I$(INCLUDE_DIR) -Wall -Wextra -g
AR := ar rcs

SRC_FILES := $(wildcard $(SRC_DIR)/*.c)
SRC_OBJS := $(SRC_FILES:.c=.o)

LIBRARY := libndbc.a

APP_SRCS := $(wildcard $(APPS_DIR)/*.c)
APPS := $(patsubst $(APPS_DIR)/%.c, $(APPS_DIR)/%, $(APP_SRCS))

all: $(LIBRARY) $(APPS)

$(LIBRARY): $(SRC_OBJS)
	$(AR) $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

$(APPS_DIR)/%: $(APPS_DIR)/%.c $(LIBRARY)
	$(CC) $(CFLAGS) $< -L. -lndbc -o $@

clean:
	rm -f $(SRC_OBJS) $(LIBRARY) $(APPS)

.PHONY: all clean
