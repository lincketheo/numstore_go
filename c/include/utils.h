#pragma once

#include <stdlib.h>

void* malloc_or_abort(size_t bytes);

void* calloc_or_abort(size_t n, size_t size);

int open_or_abort(const char* fname, int oflag, ...);

void close_wrap(int fd);

void mkdir_or_abort(const char* dirname, mode_t mode);

void rm_rf(const char* path);
