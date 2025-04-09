@startuml Configuration Process Sequence Diagram

actor Client
participant "Config" as config
participant "Driver Registry" as registry
participant "ConsoleDriver" as consoleDriver
participant "JSONFileDriver" as jsonDriver
participant "Logger" as logger

== Configuration Loading ==

Client -> config : LoadFromFile("config.yaml") or LoadFromFile("config.json")
activate config
config -> config : Read and parse YAML/JSON config
note right: Supports both YAML and JSON formats

config --> Client : return Config object
deactivate config

== Logger Creation from Config ==

Client -> config : CreateLogger()
activate config

config -> registry : Create("console", options)
activate registry
registry -> consoleDriver : NewConsoleDriver(options)
activate consoleDriver
consoleDriver --> registry : return console driver
deactivate consoleDriver
registry --> config : return console driver
deactivate registry

config -> registry : Create("json_file", options)
activate registry
registry -> jsonDriver : NewJSONFileDriver(options)
activate jsonDriver
jsonDriver --> registry : return JSON file driver
deactivate jsonDriver
registry --> config : return JSON file driver
deactivate registry

config -> logger : NewLogger(consoleDriver, jsonDriver)
activate logger
logger --> config : return logger
deactivate logger

config --> Client : return logger
deactivate config

== Logging with Configured Logger ==

Client -> logger : Info("Application started")
activate logger
note right: Uses the configured drivers
logger --> Client : success
deactivate logger

@enduml
