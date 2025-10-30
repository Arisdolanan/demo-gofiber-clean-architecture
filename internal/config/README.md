# Enhanced Logging System

## Overview

The enhanced logging system provides flexible logging capabilities with support for both Logrus and GoFiber loggers, automatic file separation by log level, and simple API usage.

```go
// Logging configuration
const USE_LOGRUS bool = true                    // Use Logrus (true) or GoFiber logger (false)
const LOG_TO_CONSOLE bool = true                // Enable console output
const LOG_TO_FILE bool = true                   // Enable file output
const LOG_SEPARATE_BY_LEVEL bool = true         // Separate files by log level
const LOG_CUSTOM_INFO_FILE string = ""          // Custom info log filename (empty = default)
const LOG_CUSTOM_WARNING_FILE string = ""       // Custom warning log filename
const LOG_CUSTOM_ERROR_FILE string = ""         // Custom error log filename
const LOG_DIRECTORY string = "storage/logs"     // Log directory path
```

## Simple API Usage

### Basic Usage (No Parameters)

```go
// GoFiber logger with enhanced features (uses configuration from constants)
app.Use(logger.New(utils.LogFiber()))

// Development/Error logger with enhanced features
app.Use(logger.New(utils.LogDev()))
```

### Custom Message Format

```go
// GoFiber logger with custom format
app.Use(logger.New(utils.LogFiber("[CUSTOM]: ${method} ${path} - ${status}")))

// Development logger with custom format  
app.Use(logger.New(utils.LogDev("[DEV]: ${method} ${path} - ${status}")))
```

## Log File Structure

When `LOG_SEPARATE_BY_LEVEL = true`:
```
storage/logs/
├── info.log        # Info, Debug, Trace logs
├── warning.log     # Warning logs
└── error.log       # Error, Fatal, Panic logs
```

When `LOG_SEPARATE_BY_LEVEL = false`:
```
storage/logs/
└── app.log         # All logs in single file
```

## Custom File Names Example

Set in constants:
```go
const LOG_CUSTOM_INFO_FILE string = "my-info.log"
const LOG_CUSTOM_WARNING_FILE string = "my-warnings.log" 
const LOG_CUSTOM_ERROR_FILE string = "my-errors.log"
```

Results in:
```
storage/logs/
├── my-info.log
├── my-warnings.log
└── my-errors.log
```

## Configuration Scenarios

### 1. Console Only Logging
```go
const LOG_TO_CONSOLE bool = true
const LOG_TO_FILE bool = false
```

### 2. File Only Logging  
```go
const LOG_TO_CONSOLE bool = false
const LOG_TO_FILE bool = true
```

### 3. Console + Single File
```go
const LOG_TO_CONSOLE bool = true
const LOG_TO_FILE bool = true
const LOG_SEPARATE_BY_LEVEL bool = false
```

### 4. Console + Level-separated Files
```go
const LOG_TO_CONSOLE bool = true
const LOG_TO_FILE bool = true
const LOG_SEPARATE_BY_LEVEL bool = true
```