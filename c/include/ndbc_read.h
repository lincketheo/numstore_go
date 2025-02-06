#pragma once

#include "ndbc_v_contig.h"

typedef struct {
  int port_num;
  v_contig fmt;
  ndbc_net_t net;
} read_args;

// Opens up port [port_num] with specified read arguments
// Ironically "read" refers to the command, so this is a
// server that strictly writes (the client is reading)
ndbc_ret_t read_server(read_args args);
