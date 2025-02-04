#include "logging.h"
#include "ndbc_dtypes.h"
#include "ndbc_read.h"
#include "ndbc_v_contig.h"
#include "ndbc_write.h"
#include "utils.h"

#include <arpa/inet.h>
#include <assert.h>
#include <complex.h>
#include <fcntl.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <unistd.h>

///////////////////////////////////////// VARIABLE DECLARATIONS

/*
 * a -> (2,5) float
 * b -> (2,) int32
 * c -> () char
 * d -> (9,10) cf64
 * e -> (9,) uint16_t
 * f -> (2, 6) uin64_t
 */

//// SHAPES
uint32_t ashape[] = { 2, 5 };
uint32_t bshape[] = { 2 };
uint32_t cshape[] = {};
uint32_t dshape[] = { 9, 10 };
uint32_t eshape[] = { 9 };
uint32_t fshape[] = { 2, 6 };

//// TYPES
#define atype float
#define btype int
#define ctype char
#define dtype complex float
#define etype uint16_t
#define ftype uint64_t

//// VALUE ARRAYS
atype avals[2000];
btype bvals[2000];
ctype cvals[2000];
dtype dvals[2000];
etype evals[2000];
ftype fvals[2000];
size_t len = 0;

//// PRINT
#define aprint(i) printf("a[%zu] = %f\n", i, avals[i])
#define bprint(i) printf("a[%zu] = %d\n", i, bvals[i])
#define cprint(i) printf("a[%zu] = %c\n", i, cvals[i])
#define dprint(i) printf("a[%zu] = %f %f\n", i, creal(dvals[i]), cimag(dvals[i]))
#define eprint(i) printf("a[%zu] = %zu\n", i, evals[i])
#define fprint(i) printf("a[%zu] = %u\n", i, fvals[i])

//// RANDOM
#define float_rand() (((float)rand() / (float)RAND_MAX) * 10)
#define char_rand() ('A' + (rand() % 26))
#define _arand() ((atype)float_rand())
#define _brand() ((btype)float_rand())
#define _crand() ((ctype)char_rand())
#define _drand() (float_rand() + I * float_rand())
#define _erand() ((etype)float_rand())
#define _frand() ((ftype)float_rand())
#define _rand(var, len0)                            \
  ({                                                \
    var##type* _head = (var##type*)head;            \
    size_t _len = (len0)*var_get_len(var);          \
    for (size_t i = 0; i < var_get_len(var); ++i) { \
      var##type temp = _##var##rand();              \
      var##vals[_len++] = temp;                     \
      *_head++ = temp;                              \
    }                                               \
    (uint8_t*)_head;                                \
  })

//// THE ACTUAL VARS
var a = {
  .shape = ashape,
  .shapel = sizeof(ashape) / sizeof *ashape,
  .dsize = sizeof(atype),
  .fd = 0,
  .ifd = 0,
};
var b = {
  .shape = bshape,
  .shapel = sizeof(bshape) / sizeof *bshape,
  .dsize = sizeof(btype),
  .fd = 0,
  .ifd = 0,
};
var c = {
  .shape = cshape,
  .shapel = sizeof(cshape) / sizeof *cshape,
  .dsize = sizeof(ctype),
  .fd = 0,
  .ifd = 0,
};
var d = {
  .shape = dshape,
  .shapel = sizeof(dshape) / sizeof *dshape,
  .dsize = sizeof(dtype),
  .fd = 0,
  .ifd = 0,
};
var e = {
  .shape = eshape,
  .shapel = sizeof(eshape) / sizeof *eshape,
  .dsize = sizeof(etype),
  .fd = 0,
  .ifd = 0,
};
var f = {
  .shape = fshape,
  .shapel = sizeof(fshape) / sizeof *fshape,
  .dsize = sizeof(ftype),
  .fd = 0,
  .ifd = 0,
};

/////////////////////////////////// V_GET_CONTIG FUNCTIONS
/// Create v_contigs

/////// DATA1
var vars11[2];
var vars12[1];
v_joined joined1[2];

/**
 * 5 x ([a, b] c)
 */
v_contig get_v_contig1()
{
  vars11[0] = a;
  vars11[1] = b;
  vars12[0] = c;

  joined1[0].vars = vars11;
  joined1[1].vars = vars12;
  joined1[0].len = 2;
  joined1[1].len = 1;

  v_contig ret = {
    .vars = joined1,
    .len = 2,
    .samples = 5,
    .i0 = 0,
  };
  return ret;
}

uint8_t* malloc_data_for_write_1(uint32_t seed, size_t* blen)
{
  srand(seed);
  int samples = 5;

  *blen = samples * var_get_blen(a);
  *blen += samples * var_get_blen(b);
  *blen += samples * var_get_blen(c);
  assert(*blen == 245);

  uint8_t* ret = malloc_or_abort(*blen);
  uint8_t* head = ret;

  for (int i = 0; i < samples; i++) {
    head = _rand(a, len + i);
    head = _rand(b, len + i);
  }

  for (int i = 0; i < samples; i++) {
    head = _rand(c, len + i);
  }
  len += samples;

  return ret;
}

/////// DATA2
var vars21[3];
var vars22[1];
var vars23[2];
v_joined joined2[3];

/**
 * 10 x ([c, b, d] e [a, f])
 */
v_contig get_v_contig2()
{
  vars21[0] = c;
  vars21[1] = b;
  vars21[2] = d;
  vars22[0] = e;
  vars23[0] = a;
  vars23[1] = f;
  joined2[0].vars = vars21;
  joined2[1].vars = vars22;
  joined2[2].vars = vars23;
  joined2[0].len = 3;
  joined2[1].len = 1;
  joined2[2].len = 2;

  v_contig ret = {
    .vars = (v_joined[]) {
        {
            .vars = (var[]) { c, b, d },
            .len = 3,
        },
        {
            .vars = (var[]) { e },
            .len = 1,
        },
        {
            .vars = (var[]) { a, f },
            .len = 1,
        },
    },
    .len = 3,
    .samples = 10,
    .i0 = 0,
  };
  return ret;
}

uint8_t* malloc_data_for_write_2(uint32_t seed, size_t* blen)
{
  srand(seed);
  int samples = 10;

  *blen = samples * var_get_blen(c);
  *blen += samples * var_get_blen(b);
  *blen += samples * var_get_blen(d);
  *blen += samples * var_get_blen(e);
  *blen += samples * var_get_blen(a);
  *blen += samples * var_get_blen(f);

  uint8_t* ret = malloc_or_abort(*blen);
  uint8_t* head = ret;

  for (int i = 0; i < samples; i++) {
    head = _rand(c, len + i);
    head = _rand(b, len + i);
    head = _rand(d, len + i);
  }

  for (int i = 0; i < samples; i++) {
    head = _rand(e, len + i);
  }

  for (int i = 0; i < samples; i++) {
    head = _rand(a, len + i);
    head = _rand(f, len + i);
  }
  len += samples;

  return ret;
}

/////// DATA3
var vars31[2];
var vars32[2];
v_joined joined3[2];

/**
 * 30 x ([c, f] [b, a])
 */
v_contig get_v_contig3()
{
  vars31[0] = c;
  vars31[1] = f;
  vars32[0] = b;
  vars32[1] = a;
  joined3[0].vars = vars31;
  joined3[1].vars = vars32;
  joined3[0].len = 2;
  joined3[1].len = 2;

  v_contig ret = {
    .vars = joined3,
    .len = 2,
    .samples = 30,
    .i0 = 0,
  };
  return ret;
}

uint8_t* malloc_data_for_write_3(uint32_t seed, size_t* blen)
{
  srand(seed);
  int samples = 30;

  *blen = samples * var_get_blen(c);
  *blen += samples * var_get_blen(f);
  *blen += samples * var_get_blen(b);
  *blen += samples * var_get_blen(a);

  uint8_t* ret = malloc_or_abort(*blen);
  uint8_t* head = ret;

  for (int i = 0; i < samples; i++) {
    head = _rand(c, len + i);
    head = _rand(f, len + i);
  }

  for (int i = 0; i < samples; i++) {
    head = _rand(b, len + i);
    head = _rand(a, len + i);
  }
  len += samples;

  return ret;
}

/////// DATA4
var vars41[1];
var vars42[2];
v_joined joined4[2];

/**
 * 3 x (b, [c, a])
 */
v_contig get_v_contig4()
{
  vars41[0] = b;
  vars42[0] = c;
  vars42[1] = a;

  joined4[0].vars = vars41;
  joined4[1].vars = vars42;
  joined4[0].len = 1;
  joined4[1].len = 2;

  v_contig ret = {
    .vars = joined4,
    .len = 2,
    .samples = 3,
    .i0 = 0,
  };
  return ret;
}

size_t blen_4(size_t samples)
{
  size_t blen = samples * var_get_blen(a);
  blen += samples * var_get_blen(b);
  blen += samples * var_get_blen(c);
  return blen;
}

uint8_t* malloc_data_for_write_4(uint32_t seed, size_t* blen)
{
  srand(seed);
  int samples = 3;

  *blen = blen_4(samples);

  uint8_t* ret = malloc_or_abort(*blen);
  uint8_t* head = ret;

  for (int i = 0; i < samples; i++) {
    head = _rand(b, len + i);
  }

  for (int i = 0; i < samples; i++) {
    head = _rand(c, len + i);
    head = _rand(a, len + i);
  }
  len += samples;

  return ret;
}

uint8_t* malloc_data_for_read_4(size_t* blen)
{
  int samples = 3;
  *blen = blen_4(samples);
  return calloc_or_abort(*blen, 1);
}

///////////////////////////////////////// OPEN
void _open()
{
  rm_rf("tests");
  mkdir_or_abort("tests", 0744);
  a.fd = open_or_abort("tests/a.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  b.fd = open_or_abort("tests/b.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  c.fd = open_or_abort("tests/c.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  d.fd = open_or_abort("tests/d.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  e.fd = open_or_abort("tests/e.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  f.fd = open_or_abort("tests/f.fd", O_RDWR | O_CREAT | O_TRUNC, 0644);

  a.ifd = open_or_abort("tests/a.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  b.ifd = open_or_abort("tests/b.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  c.ifd = open_or_abort("tests/c.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  d.ifd = open_or_abort("tests/d.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  e.ifd = open_or_abort("tests/e.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
  f.ifd = open_or_abort("tests/f.ifd", O_RDWR | O_CREAT | O_TRUNC, 0644);
}

///////////////////////////////////////// CLOSE ALL FDs
void _rewind()
{
  lseek(a.fd, 0, SEEK_SET);
  lseek(b.fd, 0, SEEK_SET);
  lseek(c.fd, 0, SEEK_SET);
  lseek(d.fd, 0, SEEK_SET);
  lseek(e.fd, 0, SEEK_SET);
  lseek(f.fd, 0, SEEK_SET);

  lseek(a.ifd, 0, SEEK_SET);
  lseek(b.ifd, 0, SEEK_SET);
  lseek(c.ifd, 0, SEEK_SET);
  lseek(d.ifd, 0, SEEK_SET);
  lseek(e.ifd, 0, SEEK_SET);
  lseek(f.ifd, 0, SEEK_SET);
}

///////////////////////////////////////// CLOSE ALL FDs
void _close()
{
  close_wrap(a.fd);
  close_wrap(b.fd);
  close_wrap(c.fd);
  close_wrap(d.fd);
  close_wrap(e.fd);
  close_wrap(f.fd);

  close_wrap(a.ifd);
  close_wrap(b.ifd);
  close_wrap(c.ifd);
  close_wrap(d.ifd);
  close_wrap(e.ifd);
  close_wrap(f.ifd);
}

///////////////////////////////////////// SERVER THREAD
void* server_write_thread(void* arg)
{
  write_args* wargs = arg;
  int* ret = malloc(sizeof *ret);
  *ret = write_server(*wargs);
  return ret;
}

void* server_read_thread(void* arg)
{
  read_args* rargs = arg;
  int* ret = malloc(sizeof *ret);
  *ret = read_server(*rargs);
  return ret;
}

///////////////////////////////////////// CLIENT ACTION
void client_write_data(uint8_t* data, size_t dlen, int port)
{
  int sock;
  struct sockaddr_in server_addr;

  // CREATE
  if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    perror("socket");
    exit(-1);
  }

  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(port);
  int ret = inet_pton(AF_INET, "127.0.0.1", &server_addr.sin_addr);
  if (ret <= 0) {
    perror("inet_pton");
    exit(-1);
  }

  // CONNECT
  ret = connect(
      sock,
      (struct sockaddr*)&server_addr,
      sizeof(server_addr));

  if (ret == -1) {
    perror("connect");
    exit(-1);
  }

  // WRITE
  if (write(sock, data, dlen) == -1) {
    perror("send");
    exit(-1);
  }

  // CLOSE
  close(sock);
}

void client_read_data(uint8_t* dest, size_t dlen, int port)
{
  int sock;
  struct sockaddr_in server_addr;

  // CREATE
  if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    perror("socket");
    exit(-1);
  }

  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(port);
  int ret = inet_pton(AF_INET, "127.0.0.1", &server_addr.sin_addr);
  if (ret <= 0) {
    perror("inet_pton");
    exit(-1);
  }

  // CONNECT
  ret = connect(
      sock,
      (struct sockaddr*)&server_addr,
      sizeof(server_addr));

  if (ret == -1) {
    perror("connect");
    exit(-1);
  }

  // READ
  if (read(sock, dest, dlen) == -1) {
    perror("send");
    exit(-1);
  }

  // CLOSE
  close(sock);
}

int main()
{
  // Return value from threads
  int* _ret;
  int ret = 0;

  // First, open all files
  _open();

  // Allocate memory buffers
  size_t wdlen1, rdlen4;
  uint8_t* wdata1 = malloc_data_for_write_1(1234, &wdlen1);
  uint8_t* rdata4 = malloc_data_for_read_4(&rdlen4);

  // Create the server write thread
  pthread_t thread_id;
  write_args args = {
    .fmt = get_v_contig1(),
    .net = TCP,
    .port_num = 8080,
  };
  if (pthread_create(&thread_id, NULL, server_write_thread, &args) != 0) {
    perror("pthread_create");
    return -1;
  }

  // Sleep for a bit - no synchro
  usleep(100);

  // Write to the port
  client_write_data(wdata1, wdlen1, 8080);

  // Join the server thread
  pthread_join(thread_id, (void**)&_ret);
  ret = ret | *_ret;

  // Rewind all file pointers to the start
  _rewind();

  // Create the server read thread
  read_args rargs = {
    .fmt = get_v_contig4(),
    .net = TCP,
    .port_num = 8081,
  };
  if (pthread_create(&thread_id, NULL, server_read_thread, &rargs) != 0) {
    perror("pthread_create");
    return -1;
  }

  // Sleep for a bit - no synchro
  usleep(1000);

  // Read from the port
  client_read_data(rdata4, rdlen4, 8081);

  // Join the server thread
  pthread_join(thread_id, (void**)&_ret);
  ret = ret | *_ret;

  // Free data
  free(wdata1);
  free(rdata4);
  free(_ret);

  // Close all files
  _close();

  return ret;
}
