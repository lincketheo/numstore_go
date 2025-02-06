#pragma once

#include <stdint.h>
#include <stdlib.h>

typedef enum {
  SUCCESS,
  FAILURE,
} ndbc_ret_t;

typedef enum {
  UDP,
  TCP,
} ndbc_net_t;

typedef struct {
  enum {
    INDEX_INT,
    INDEX_SLICE,
  } type;
  union {
    int index;
    struct {
      int start;
      int stop;
      int step;
    };
  };
} index_t;

typedef struct {
  uint32_t *shape; // The data type shape
  size_t shapel;   // The shape length

  size_t dsize; // The size of each element
  int fd;       // The raw data file
  int ifd;      // The file pointer of the index

  // Meaningless for write
  // index_t* indexes; // Indexes a[1,:,1:2] -> (1, :, 1:2)
  // size_t indlen;    // len(indexes)
} var;

size_t var_get_len(var a);
size_t var_get_blen(var a);
