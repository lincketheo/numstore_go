#include "ndbc_read.h"
#include "logging.h"
#include "ndbc_dtypes.h"
#include "ndbc_v_contig.h"

#include <assert.h>
#include <netinet/in.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <unistd.h>

static int handle_read_client(v_contig fmt, int client_fd); // Callback
static void cleanup(int signum);                            // Signal Quit

ndbc_ret_t read_server(read_args args) {
  log_debug("Starting read server on port: %d\n", args.port_num);

  // PREDECL
  int server_fd, client_fd;
  struct sockaddr_in server_addr, client_addr;
  socklen_t client_len = sizeof(client_addr);

  signal(SIGINT, cleanup);

  // CREATE
  if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
    perror("socket");
    return FAILURE;
  }

  server_addr.sin_family = AF_INET;
  server_addr.sin_addr.s_addr = INADDR_ANY;
  server_addr.sin_port = htons(args.port_num);

  // BIND
  int ret =
      bind(server_fd, (struct sockaddr *)&server_addr, sizeof(server_addr));
  if (ret == -1) {
    perror("bind");
    close(server_fd);
    return FAILURE;
  }

  // LISTEN
  ret = listen(server_fd, 5);
  if (ret == -1) {
    perror("listen");
    close(server_fd);
    return FAILURE;
  }
  log_info("Read server listening on port: %d\n", args.port_num);

  // ACCEPT
  client_fd = accept(server_fd, (struct sockaddr *)&client_addr, &client_len);
  if (client_fd == -1) {
    perror("accept");
    goto failed_loop;
  }

  // HANDLE
  log_info("New Connection\n");
  if (handle_read_client(args.fmt, client_fd)) {
    goto failed_loop;
  }

  // CLOSE
  if (close(client_fd) == -1) {
    perror("close");
    goto failed_loop;
  }

failed_loop:
  return EXIT_FAILURE;
}

int handle_read_client(v_contig fmt, int client_fd) {
  int ret = 0;

  // ALLOC
  v_contig_mem_space s = {0};
  if (v_contig_mem_space_alloc(&s, fmt)) {
    ret = -1;
    goto theend;
  }

  // READ
  log_debug("Reading from files\n");
  if (v_contig_read(&s, fmt) == -1) {
    ret = -1;
    goto theend;
  }

  // WRITE
  log_debug("Writing %zu bytes to socket\n", s.blen);
  ssize_t nwrite = write(client_fd, s.raveled, s.blen);
  if (nwrite == -1 || (size_t)nwrite != s.blen) { // TODO - handle err
    perror("write");
    ret = -1;
    goto theend;
  }

  // DEALLOC
  if (v_contig_mem_space_free(&s)) {
    ret = -1;
    goto theend;
  }

theend:
  return ret;
}

void cleanup(int signum) {
  printf("Shutting down write server from signal: %d\n", signum);
  exit(0);
}
