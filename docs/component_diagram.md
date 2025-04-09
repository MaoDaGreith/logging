@startuml Logging Package Component Diagram

package "Application" {
  [Client Code] as client
}

package "Logging Library" {
  package "Core" {
    [Logger] as logger
    [Transaction] as transaction
    [Log Entry] as logEntry
  }

  package "Drivers" {
    [Driver Registry] as registry
    [Console Driver] as consoleDriver
    [JSON File Driver] as jsonDriver
    [Text File Driver] as textDriver
    
    registry --> consoleDriver : creates
    registry --> jsonDriver : creates
    registry --> textDriver : creates
  }

  package "Configuration" {
    [Config Loader] as configLoader
    
    configLoader --> registry : uses
    configLoader --> logger : creates
  }
  
  logger --> transaction : creates
  logger --> logEntry : creates
  transaction --> logEntry : creates
  logger --> consoleDriver : uses
  logger --> jsonDriver : uses
  logger --> textDriver : uses
}

package "Output Destinations" {
  [Console (stdout/stderr)] as console
  [JSON Log File] as jsonFile
  [Text Log File] as textFile
}

client --> logger : uses
client --> transaction : uses
client --> configLoader : uses

consoleDriver --> console : writes to
jsonDriver --> jsonFile : writes to
textDriver --> textFile : writes to

note right of [Config Loader]
  Configuration supports:
  - YAML format
  - JSON format
  - Environment variables
  - Default values
end note

@enduml
``` 