#pragma once

#include <stdarg.h>
#include <stdio.h>

#define DEBUG 0
#define INFO 1
#define WARN 2
#define ERROR 3

#ifndef LOG_LEVEL
#define LOG_LEVEL DEBUG // Default log level
#endif

#ifndef LOG_FILE
#define LOG_FILE stderr // Default log output
#endif

void log_debug(const char *fmt, ...);

void log_info(const char *fmt, ...);

void log_warn(const char *fmt, ...);

void log_error(const char *fmt, ...);
