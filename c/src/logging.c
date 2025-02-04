#include "logging.h"

static inline void log_message(
    int level,
    const char* level_str,
    const char* fmt,
    va_list args)
{
  if (level >= LOG_LEVEL) {
    FILE* file = LOG_FILE;
    fprintf(file, "[%-5s]       ", level_str);
    vfprintf(file, fmt, args);
    va_end(args);
  }
}

void log_debug(const char* fmt, ...)
{
  if (DEBUG >= LOG_LEVEL) {
    va_list args;
    va_start(args, fmt);
    log_message(DEBUG, "DEBUG", fmt, args);
    va_end(args);
  }
}

void log_info(const char* fmt, ...)
{
  if (INFO >= LOG_LEVEL) {
    va_list args;
    va_start(args, fmt);
    log_message(INFO, "INFO", fmt, args);
    va_end(args);
  }
}

void log_warn(const char* fmt, ...)
{
  if (WARN >= LOG_LEVEL) {
    va_list args;
    va_start(args, fmt);
    log_message(WARN, "WARN", fmt, args);
    va_end(args);
  }
}

void log_error(const char* fmt, ...)
{
  if (ERROR >= LOG_LEVEL) {
    va_list args;
    va_start(args, fmt);
    log_message(ERROR, "ERROR", fmt, args);
    va_end(args);
  }
}
