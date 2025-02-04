#pragma once

#include "ndbc_v_contig.h"

/**
 * If UDP, it requires minimum number
 * of data for 1 transaction to be less
 * than maximum UDP datagram length
 */
typedef struct {
  int port_num;
  v_contig fmt;
  ndbc_net_t net;
} write_args;

// Opens up port [port_num] with specified write arguments
// Ironically "write" refers to the command, so this is a
// server that strictly reads (the client is writing to me)
ndbc_ret_t
write_server(write_args args);
