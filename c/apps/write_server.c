#include "ndbc_dtypes.h"
#include "ndbc_write.h"

#include <assert.h>
#include <fcntl.h>
#include <stdio.h>
#include <unistd.h>

int afd;
int bfd;
int cfd;
int aifd;
int bifd;
int cifd;

uint32_t ashape[] = { 2, 5 };
uint32_t bshape[] = { 2 };
uint32_t cshape[] = {};

var vars1[] = {
  {
      .shape = ashape,
      .shapel = 2,
      .dsize = 4,
      .fd = 0,
  },
  {
      .shape = bshape,
      .shapel = 1,
      .dsize = 4,
      .fd = 0,
  },
};

var vars2[] = {
  {
      .shape = cshape,
      .shapel = 0,
      .dsize = 1,
      .fd = 0,
  },
};

v_joined joined[] = {
  {
      .vars = vars1,
      .len = 2,
  },
  {
      .vars = vars2,
      .len = 1,
  }
};

write_args args = {
  .port_num = 8080,
  .fmt = {
      .vars = joined,
      .len = 2,
      .samples = 5,
      .i0 = 100,
  },
  .net = TCP,
};

void _open()
{
  afd = open("a", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  bfd = open("b", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  cfd = open("c", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  assert(afd != -1);
  assert(bfd != -1);
  assert(cfd != -1);

  aifd = open("ai", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  bifd = open("bi", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  cifd = open("ci", O_WRONLY | O_CREAT | O_TRUNC, 0644);
  assert(aifd != -1);
  assert(bifd != -1);
  assert(cifd != -1);
}

void _close()
{
  close(afd);
  close(bfd);
  close(cfd);
  close(aifd);
  close(bifd);
  close(cifd);
}

static write_args get_write_args()
{
  args.fmt.vars[0].vars[0].fd = afd;
  args.fmt.vars[0].vars[1].fd = bfd;
  args.fmt.vars[1].vars[0].fd = cfd;
  args.fmt.vars[0].vars[0].ifd = aifd;
  args.fmt.vars[0].vars[1].ifd = bifd;
  args.fmt.vars[1].vars[0].ifd = cifd;
  return args;
}

int main()
{
  _open();
  int ret = write_server(get_write_args());
  _close();
  return ret;
}
