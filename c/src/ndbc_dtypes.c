#include "ndbc_dtypes.h"

#include <assert.h>

size_t var_get_blen(var a) { return var_get_len(a) * a.dsize; }

size_t var_get_len(var a) {
  size_t ret = 1;
  for (size_t i = 0; i < a.shapel; ++i) {
    assert(a.shape[i] != 0);
    ret *= a.shape[i];
  }
  return ret;
}
