@startuml Logging Process Sequence Diagram

actor Client
participant "Logger" as logger
participant "Transaction" as transaction
participant "LogEntry" as logEntry
participant "ConsoleDriver" as consoleDriver
participant "JSONFileDriver" as jsonDriver
participant "TextFileDriver" as textDriver
participant "stdout/stderr" as console
participant "JSON File" as jsonFile
participant "Text File" as textFile

== Basic Logging ==

Client -> logger : Info("User login", attributes)
activate logger

logger -> logEntry : create(timestamp, Info, "User login", attributes, "")
activate logEntry
logEntry --> logger : return entry
deactivate logEntry

logger -> consoleDriver : Log(entry)
activate consoleDriver
consoleDriver -> consoleDriver : format(entry)
consoleDriver -> console : write formatted log
console --> consoleDriver : success
consoleDriver --> logger : success
deactivate consoleDriver

logger -> jsonDriver : Log(entry)
activate jsonDriver
jsonDriver -> jsonDriver : check minLevel
jsonDriver -> jsonFile : write JSON entry
jsonFile --> jsonDriver : success
jsonDriver --> logger : success
deactivate jsonDriver

logger -> textDriver : Log(entry)
activate textDriver
textDriver -> textDriver : formatLogEntry(entry)
textDriver -> textFile : write formatted log
textFile --> textDriver : success
textDriver --> logger : success
deactivate textDriver

logger --> Client : success
deactivate logger

== Transaction-based Logging ==

Client -> logger : NewTransaction("request-123")
activate logger
logger -> transaction : create("request-123", logger)
activate transaction
transaction --> logger : return transaction
logger --> Client : transaction
deactivate logger

Client -> transaction : Info("Processing request")
activate transaction
transaction -> logEntry : create(timestamp, Info, "Processing request", {}, "request-123")
activate logEntry
logEntry --> transaction : return entry
deactivate logEntry

transaction -> consoleDriver : Log(entry)
activate consoleDriver
consoleDriver -> consoleDriver : format(entry)
consoleDriver -> console : write formatted log with transaction ID
console --> consoleDriver : success
consoleDriver --> transaction : success
deactivate consoleDriver

transaction -> jsonDriver : Log(entry)
activate jsonDriver
jsonDriver -> jsonDriver : check minLevel
jsonDriver -> jsonFile : write JSON entry with transaction ID
jsonFile --> jsonDriver : success
jsonDriver --> transaction : success
deactivate jsonDriver

transaction -> textDriver : Log(entry)
activate textDriver
textDriver -> textDriver : formatLogEntry(entry)
textDriver -> textFile : write formatted log with transaction ID
textFile --> textDriver : success
textDriver --> transaction : success
deactivate textDriver

transaction --> Client : success
deactivate transaction

== Logger Closure ==

Client -> logger : Close()
activate logger

logger -> consoleDriver : Close()
activate consoleDriver
consoleDriver --> logger : success (no-op)
deactivate consoleDriver

logger -> jsonDriver : Close()
activate jsonDriver
jsonDriver -> jsonFile : flush and close
jsonFile --> jsonDriver : success
jsonDriver --> logger : success
deactivate jsonDriver

logger -> textDriver : Close()
activate textDriver
textDriver -> textFile : flush and close
textFile --> textDriver : success
textDriver --> logger : success
deactivate textDriver

logger --> Client : success
deactivate logger

@enduml 
