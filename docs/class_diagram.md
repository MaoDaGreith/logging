@startuml
title Logging Package Class Diagram

package "core" {
  enum Level {
    Debug
    Info
    Warning
    Error
    +String(): string
  }

  class LogEntry {
    +Timestamp: time.Time
    +Level: Level
    +Message: string
    +Attrs: Attributes
    +TransactionID: string
  }

  interface Logger {
    +Debug(msg string, attrs ...Attributes): error
    +Info(msg string, attrs ...Attributes): error
    +Warning(msg string, attrs ...Attributes): error
    +Error(msg string, attrs ...Attributes): error
    +Log(level Level, msg string, attrs ...Attributes): error
    +NewTransaction(txID string): Transaction
    +Close(): error
  }

  class logger {
    -drivers: []Driver
    +Debug(msg string, attrs ...Attributes): error
    +Info(msg string, attrs ...Attributes): error
    +Warning(msg string, attrs ...Attributes): error
    +Error(msg string, attrs ...Attributes): error
    +Log(level Level, msg string, attrs ...Attributes): error
    +NewTransaction(txID string): Transaction
    +Close(): error
  }

  interface Transaction {
    +Debug(msg string, attrs ...Attributes): error
    +Info(msg string, attrs ...Attributes): error
    +Warning(msg string, attrs ...Attributes): error
    +Error(msg string, attrs ...Attributes): error
    +Log(level Level, msg string, attrs ...Attributes): error
    +ID(): string
  }

  class transaction {
    -id: string
    -logger: *logger
    +Debug(msg string, attrs ...Attributes): error
    +Info(msg string, attrs ...Attributes): error
    +Warning(msg string, attrs ...Attributes): error
    +Error(msg string, attrs ...Attributes): error
    +Log(level Level, msg string, attrs ...Attributes): error
    +ID(): string
  }

  interface Driver {
    +Log(entry *LogEntry): error
    +Close(): error
  }

  class Attributes {
    Map[string]string
  }

  Logger <|.. logger
  Transaction <|.. transaction
  logger ..> Driver : uses
  transaction --> logger : references
  logger ..> LogEntry : creates
  transaction ..> LogEntry : creates
  LogEntry *-- Level
  LogEntry *-- Attributes
}

package "drivers" {
  class ConsoleDriver {
    -stdout: io.Writer
    -stderr: io.Writer
    -minLevel: core.Level
    -timeFormat: string
    -colorized: bool
    +Log(entry *core.LogEntry): error
    +Close(): error
    -format(entry *core.LogEntry): string
    -colorizeLevel(level core.Level, levelStr string): string
  }

  class JSONFileDriver {
    -filePath: string
    -file: *os.File
    -encoder: *json.Encoder
    -minLevel: core.Level
    -mu: sync.Mutex
    +Log(entry *core.LogEntry): error
    +Close(): error
  }

  class TextFileDriver {
    -filePath: string
    -file: *os.File
    -minLevel: core.Level
    -timeFormat: string
    -mu: sync.Mutex
    +Log(entry *core.LogEntry): error
    +Close(): error
    -formatLogEntry(entry *core.LogEntry): string
  }

  class DriverRegistry {
    -registry: map[string]DriverConstructor
    +Register(name string, constructor DriverConstructor)
    +Create(name string, options map[string]interface{}): (core.Driver, error)
  }

  core.Driver <|.. ConsoleDriver
  core.Driver <|.. JSONFileDriver
  core.Driver <|.. TextFileDriver
}

package "config" {
  class Config {
    +DefaultLevel: string
    +Drivers: []DriverConfig
    +LoadFromFile(path string): (*Config, error)
    +LoadDefault(): (*Config, error)
    +CreateLogger(): (core.Logger, error)
    +SaveToFile(path string): error
  }

  class DriverConfig {
    +Type: string
    +MinLevel: string
    +Options: map[string]interface{}
  }

  Config *-- DriverConfig
  Config ..> drivers.DriverRegistry : uses
  Config ..> core.Logger : creates
}

drivers.ConsoleDriver ..> core.LogEntry : processes
drivers.JSONFileDriver ..> core.LogEntry : processes
drivers.TextFileDriver ..> core.LogEntry : processes
drivers.DriverRegistry ..> core.Driver : creates

@enduml