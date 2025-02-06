#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

void error_exit(const char *msg) {
  perror(msg);
  exit(EXIT_FAILURE);
}

int main(int argc, char *argv[]) {
  if (argc != 2) {
    fprintf(stderr, "Usage: %s <PORT>\n", argv[0]);
    return EXIT_FAILURE;
  }

  // PARSE ARGS
  int port = atoi(argv[1]);
  if (port <= 0 || port > 65535) {
    fprintf(stderr, "Invalid port number\n");
    return EXIT_FAILURE;
  }

  int sock;
  struct sockaddr_in server_addr;
  char buffer[2048];

  // CREATE
  if ((sock = socket(AF_INET, SOCK_STREAM, 0)) == -1)
    error_exit("socket");

  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(port);
  int ret = inet_pton(AF_INET, "127.0.0.1", &server_addr.sin_addr);
  if (ret <= 0) {
    error_exit("inet_pton");
  }

  // CONNECT
  ret = connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr));

  if (ret == -1) {
    error_exit("connect");
  }

  // WRITE
  ssize_t readb;
  while ((readb = fread(buffer, 1, 2048, stdin)) > 0) {
    if (send(sock, buffer, readb, 0) == -1)
      error_exit("send");
  }

  // CLOSE
  close(sock);
  return 0;
}
