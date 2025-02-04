#include "ndbc_v_contig.h"
#include "logging.h"
#include "ndbc_dtypes.h"

#include <assert.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <unistd.h>

size_t v_contig_get_num_vars(v_contig fmt)
{
  size_t ret = 0;
  for (size_t i = 0; i < fmt.len; ++i) {
    ret += fmt.vars[i].len;
  }
  return ret;
}

int v_contig_mem_space_alloc(v_contig_mem_space* dest, v_contig src)
{
  assert(dest);
  assert(!dest->raveled);
  assert(!dest->unraveled);
  assert(!dest->ibuffer);
  assert(!dest->vars);
  assert(!dest->index_mmaps);

  int ret = 0;

  dest->blen = v_contig_blen(src);

  // ALLOC
  dest->unraveled = mmap(NULL, dest->blen, PROT_READ | PROT_WRITE,
      MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
  if (dest->unraveled == NULL) {
    perror("mmap");
    ret = -1;
    goto theend;
  }

  dest->raveled = mmap(NULL, dest->blen, PROT_READ | PROT_WRITE,
      MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
  if (dest->raveled == NULL) {
    perror("mmap");
    ret = -1;
    munmap(dest->unraveled, dest->blen);
    goto theend;
  }

  dest->samples = src.samples;
  dest->ibuflen = 0;
  dest->i0 = src.i0;
  dest->ibuffer = malloc(dest->samples * sizeof *dest->ibuffer);
  if (dest->ibuffer == NULL) {
    perror("malloc");
    ret = -1;
    munmap(dest->unraveled, dest->blen);
    munmap(dest->raveled, dest->blen);
    goto theend;
  }

  dest->num_vars = v_contig_get_num_vars(src);
  dest->vars = malloc(dest->num_vars * sizeof *dest->vars);
  if (dest->vars == NULL) {
    perror("malloc");
    ret = -1;
    munmap(dest->unraveled, dest->blen);
    munmap(dest->raveled, dest->blen);
    free(dest->ibuffer);
    goto theend;
  }

  // FILL VARS
  size_t k = 0;
  for (size_t i = 0; i < src.len; ++i) {
    for (size_t j = 0; j < src.vars[i].len; ++j) {
      dest->vars[k++] = src.vars[i].vars[j];
    }
  }

theend:
  return ret;
}

int v_contig_mem_space_free(v_contig_mem_space* dest)
{
  assert(dest);
  assert(dest->raveled);
  assert(dest->unraveled);
  assert(dest->ibuffer);
  assert(dest->vars);
  assert(!dest->index_mmaps);

  int ret = 0;
  if (munmap(dest->raveled, dest->blen)) {
    perror("munmap");
    ret = -1;
  }
  if (munmap(dest->unraveled, dest->blen)) {
    perror("munmap");
    ret = -1;
  }
  free(dest->ibuffer);
  free(dest->vars);
  return ret;
}

/**
 * fmt:  "[a, b], c"
 * src:  [a, b, a, b, a, b, c, c, c]
 * dest: [a, a, a, b, b, b, c, c, c]
 *
 * Same algorithm as ravel except memcpy(a, b) is memcpy(b, a)
 */
static void unravel(uint8_t* dest, const uint8_t* src, v_contig fmt)
{
  uint8_t* contig_offset = dest; // Offset of contig[i]
  uint8_t* joined_offset = dest; // Offset of joined[i]

  // For each [a, b], c in "[a, b], c"
  for (size_t j = 0; j < fmt.len; ++j) {
    v_joined joined = fmt.vars[j];

    // For each c, in [c, c, c, c, c]
    for (size_t c = 0; c < fmt.samples; ++c) {

      joined_offset = contig_offset;

      // For each a, b in "[a, b]"
      for (size_t v = 0; v < joined.len; ++v) {
        var _var = joined.vars[v];

        // Byte size of one element of a
        size_t len = var_get_blen(_var);

        memcpy(joined_offset + c * len, src, len);
        joined_offset += len * fmt.samples;
        src += len;
      }
    }

    contig_offset = joined_offset;
  }
}

/**
 * fmt:  "[a, b], c"
 * src:  [a, a, a, b, b, b, c, c, c]
 * dest: [a, b, a, b, a, b, c, c, c]
 *
 * Same algorithm as unravel except memcpy(a, b) is memcpy(b, a)
 */
static void ravel(uint8_t* dest, const uint8_t* src, v_contig fmt)
{
  const uint8_t* contig_offset = src; // Offset of contig[i]
  const uint8_t* joined_offset = src; // Offset of joined[i]

  // For each [a, b], c in "[a, b], c"
  for (size_t j = 0; j < fmt.len; ++j) {
    v_joined joined = fmt.vars[j];

    // For each c, in [c, c, c, c, c]
    for (size_t c = 0; c < fmt.samples; ++c) {

      joined_offset = contig_offset;

      // For each a, b in "[a, b]"
      for (size_t v = 0; v < joined.len; ++v) {
        var _var = joined.vars[v];

        // Byte size of one element of a
        size_t len = var_get_blen(_var);

        memcpy(dest, joined_offset + c * len, len);
        joined_offset += len * fmt.samples;
        dest += len;
      }
    }

    contig_offset = joined_offset;
  }
}

/**
 * fmt: "[a, b], c
 * src: [a, a, a, b, b, b, c, c, c]
 * write(afd,   [a, a, a])
 * write(bfd,   [b, b, b])
 * write(cfd,   [c, c, c])
 * write(aifd,  [i0, i0 + 1, i0 + 2])
 * write(bifd,  [i0, i0 + 1, i0 + 2])
 * write(cifd,  [i0, i0 + 1, i0 + 2])
 */
static int write_unraveled(v_contig_mem_space* s)
{
  assert(s);
  assert(s->unraveled);
  assert(s->ibuffer);
  assert(s->vars);

  int ret = 0;
  uint8_t* head = s->unraveled;

  // Fill Indexes with values
  for (size_t i = 0; i < s->samples; ++i) {
    s->ibuffer[i] = s->i0 + i;
  }

  for (size_t i = 0; i < s->num_vars; ++i) {
    var _var = s->vars[i];

    // WRITE DATA
    size_t towrite = var_get_blen(_var) * s->samples;
    ssize_t nwrite = write(_var.fd, head, towrite);
    log_debug("Writing to variable: %zu\n", i);
    if (nwrite == -1 || (size_t)nwrite != towrite) {
      perror("write");
      ret = -1;
      goto theend;
    }
    head += towrite;

    // WRITE INDEXES
    towrite = s->samples * sizeof *s->ibuffer;
    nwrite = write(_var.ifd, s->ibuffer, towrite);
    log_debug("Writing indexes to variable: %zu\n", i);
    if (nwrite == -1 || (size_t)nwrite != towrite) {
      perror("write");
      ret = -1;
      goto theend;
    }
  }

theend:
  return ret;
}

static int v_contig_mem_space_munmap_indexes_permissive(
    v_contig_mem_space* s)
{
  assert(s);

  int ret = 0;

  if (s->index_mmaps) {
    assert(s->index_map_lens);
    // Call munmap on each initialized map
    for (size_t i = 0; i < s->num_vars; ++i) {
      if (s->index_mmaps[i]) {
        size_t blen = s->index_map_lens[i] * sizeof *s->index_mmaps[i];
        log_debug("Unmapping memory of indexes for variable: %zu\n", i);
        if (munmap(s->index_mmaps[i], blen)) {
          perror("munmap");
          ret = -1;
        }
      }
    }
    free(s->index_mmaps);
    s->index_mmaps = NULL;
  }
  if (s->index_map_lens) {
    free(s->index_map_lens);
    s->index_map_lens = NULL;
  }

  return ret;
}

static int v_contig_mem_space_mmap_indexes(v_contig_mem_space* s)
{
  assert(s);
  assert(!s->index_mmaps);
  assert(!s->index_map_lens);

  s->index_mmaps = calloc(s->num_vars, sizeof *s->index_mmaps);
  if (s->index_mmaps == NULL) {
    goto failed;
  }
  s->index_map_lens = malloc(s->num_vars * sizeof *s->index_map_lens);
  if (s->index_map_lens == NULL) {
    goto failed;
  }

  // Call mmap for each index file
  for (size_t i = 0; i < s->num_vars; ++i) {
    // Get file size
    struct stat st;
    if (fstat(s->vars[i].ifd, &st) == -1) {
      perror("fstat");
      goto failed;
    }

    size_t blen = st.st_size;
    s->index_map_lens[i] = blen / sizeof *s->index_mmaps[i];

    log_debug("Mapping memory %zu bytes for variable: %zu fd: %d\n",
        blen, i, s->vars[i].ifd);
    s->index_mmaps[i] = mmap(NULL, blen,
        PROT_READ | PROT_WRITE, MAP_PRIVATE,
        s->vars[i].ifd, 0);
    if (s->index_mmaps[i] == MAP_FAILED) {
      perror("mmap");
      goto failed;
    }
  }

  return 0;

failed:
  // Ignore return value - we already failed
  v_contig_mem_space_munmap_indexes_permissive(s);
  return -1;
}

/**
 * Finds shared indexes in index arrays with limit
 */
static int index_find_shared(
    size_t* dest,     // output array
    size_t* dlen,     // final count
    size_t** indexes, // list of sorted arrays
    size_t* ilens,    // length of each array
    size_t len,       // number of arrays
    size_t limit)     // maximum results
{
  size_t count = 0;
  size_t pos[64]; // or dynamically allocated if len > 64
  if (len > 64) {
    return -1;
  }

  // Initialize positions
  for (size_t i = 0; i < len; i++) {
    pos[i] = 0;
  }

  while (1) {
    // Find the max index
    // at the current pos slice
    size_t max_val = 0;
    for (size_t i = 0; i < len; i++) {

      // Reached the end of one array - no more
      // common intersections possible
      if (pos[i] >= ilens[i]) {
        goto done;
      }

      //
      size_t val = indexes[i][pos[i]];
      if (val > max_val) {
        max_val = val;
      }
    }

    int matched = 1; // whether or not we should include max_val

    // Advance all indexes up to max_val
    for (size_t i = 0; i < len; i++) {
      while (pos[i] < ilens[i] && indexes[i][pos[i]] < max_val) {
        pos[i]++;
      }

      // Reached the end of one array - no more
      // common intersections possible
      if (pos[i] >= ilens[i]) {
        goto done;
      }
      if (indexes[i][pos[i]] != max_val) {
        matched = 0;
      }
    }

    // If all matched, store
    if (matched) {
      dest[count++] = max_val;
      if (count == limit) {
        break;
      }

      // Increment all positions by 1
      for (size_t i = 0; i < len; i++) {
        pos[i]++;
        if (pos[i] >= ilens[i]) {
          goto done;
        }
      }
    }
  }

done:
  *dlen = count;
  return 0;
}

static ssize_t argfind(size_t* indexes, size_t ilen, size_t find)
{
  size_t left = 0, right = ilen;

  while (left < right) {
    size_t mid = left + (right - left) / 2;
    if (indexes[mid] == find) {
      return mid;
    } else if (indexes[mid] < find) {
      left = mid + 1;
    } else {
      right = mid;
    }
  }

  return -1;
}

static int read_unraveled(v_contig_mem_space* s)
{
  assert(s);
  assert(s->unraveled);
  assert(s->ibuffer);
  assert(s->vars);

  int ret = 0;

  // Memory map index files
  if (v_contig_mem_space_mmap_indexes(s)) {
    ret = -1;
    goto theend;
  }

  // Fill size buffer with shared indexes
  if (index_find_shared(
          s->ibuffer,
          &s->ibuflen,
          s->index_mmaps,
          s->index_map_lens,
          s->num_vars,
          s->samples)) {
    ret = -1;
    goto theend;
  }

  if (LOG_LEVEL > DEBUG) {
    log_debug("Shared indexes for read are:\n");
    for (size_t i = 0; i < s->ibuflen; ++i) {
      log_debug("%d\n", s->ibuffer[i]);
    }
  }

  uint8_t* head = s->unraveled;
  for (size_t i = 0; i < s->num_vars; ++i) {
    for (size_t j = 0; j < s->ibuflen; ++j) {
      var _var = s->vars[i];
      size_t blen = var_get_blen(_var);

      // FIND OFFSET
      size_t index = s->ibuffer[j];
      ssize_t var_index = argfind(
          s->index_mmaps[i],
          s->index_map_lens[i],
          index);
      assert(var_index != -1);

      // READ
      ssize_t nread = read(_var.fd, head, blen);
      if (nread == -1 || (size_t)nread != blen) {
        perror("read");
        ret = -1;
        goto theend;
      }
      head += blen;
    }
  }

theend:
  // Memory unmap index files
  if (v_contig_mem_space_munmap_indexes_permissive(s)) {
    ret = -1;
  }
  return ret;
}

int v_contig_write(v_contig_mem_space* s, v_contig fmt)
{
  int ret = 0;

  // UNRAVEL
  unravel(s->unraveled, s->raveled, fmt);

  // WRITE
  if ((ret = write_unraveled(s))) {
    return ret;
  }

  return 0;
}

int v_contig_read(v_contig_mem_space* s, v_contig fmt)
{
  int ret = 0;

  // READ
  if ((ret = read_unraveled(s))) {
    return ret;
  }

  // RAVEL
  ravel(s->raveled, s->unraveled, fmt);

  return 0;
}

size_t v_joined_blen(v_joined v)
{
  size_t ret = 0;
  for (size_t i = 0; i < v.len; ++i) {
    ret += var_get_blen(v.vars[i]);
  }
  return ret;
}

size_t v_contig_blen(v_contig v)
{
  size_t ret = 0;
  for (size_t i = 0; i < v.len; ++i) {
    ret += v_joined_blen(v.vars[i]);
  }
  return ret * v.samples;
}
