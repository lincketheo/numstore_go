#include "utils.h"

#include <fcntl.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <unistd.h>

void *malloc_or_abort(size_t bytes) {
  void *ret = malloc(bytes);
  if (ret == NULL) {
    perror("malloc");
    exit(-1);
  }
  return ret;
}

void *calloc_or_abort(size_t n, size_t size) {
  void *ret = calloc(n, size);
  if (ret == NULL) {
    perror("malloc");
    exit(-1);
  }
  return ret;
}

int open_or_abort(const char *fname, int oflag, ...) {
  int fd;

  if (oflag & O_CREAT) {
    va_list args;
    va_start(args, oflag);
    mode_t mode = va_arg(args, mode_t);
    va_end(args);

    fd = open(fname, oflag, mode);
  } else {
    fd = open(fname, oflag);
  }

  if (fd == -1) {
    perror("open_or_abort");
    exit(EXIT_FAILURE);
  }

  return fd;
}

void close_wrap(int fd) {
  if (close(fd)) {
    perror("close");
  }
}

void mkdir_or_abort(const char *dirname, mode_t mode) {
  if (mkdir(dirname, mode)) {
    perror("mkdir");
    exit(EXIT_FAILURE);
  }
}

void rm_rf(const char *path) {
  char cmd[1024];
  snprintf(cmd, sizeof(cmd), "rm -rf -- '%s'", path);
  system(cmd);
}
