#pragma once

#include "ndbc_dtypes.h"

/*
 * Joined vars
 * len = 3
 * [a, b, c] => [a, b, c, a, b, c, a, b, c]
 */
typedef struct {
  var* vars;
  size_t len; // len(vars)
} v_joined;

/**
 * Contiguous vars
 * len = 2
 * samples = 3
 * a, [b, c] => [a, a, a, b, c, b, c, b, c]
 */
typedef struct {
  v_joined* vars;
  size_t len;     // len(vars)
  size_t samples; // Num elements for each contiguous packet
  size_t i0;
} v_contig;

// Returns the number of variables in v_contig
size_t v_contig_get_num_vars(v_contig);

size_t v_joined_blen(v_joined v);

// Returns total bytes of v_contig packet
size_t v_contig_blen(v_contig v);

typedef struct {
  // Temporary buffers
  uint8_t* raveled;
  uint8_t* unraveled;
  size_t blen;

  // Buffer allocated for indexes
  size_t* ibuffer;
  size_t ibuflen; // Actual length of ibuffer
  size_t samples; // Desired length of ibuffer
  size_t i0;

  // All vars - ordered in occurrence, without join/contig structure
  var* vars;
  size_t num_vars;

  // Temporary array of index mmaps - should be opened
  // / closed temporarily - not kept open the whole time
  size_t** index_mmaps;
  size_t* index_map_lens;
} v_contig_mem_space;

// RAVEL/UNRAVEL alloc / free
int v_contig_mem_space_alloc(
    v_contig_mem_space* dest,
    v_contig src);

int v_contig_mem_space_free(
    v_contig_mem_space* dest);

// WRITE out of
int v_contig_write(
    v_contig_mem_space* s,
    v_contig fmt);

// READ into
int v_contig_read(
    v_contig_mem_space* s,
    v_contig fmt);
