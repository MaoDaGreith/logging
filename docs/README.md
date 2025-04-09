# UML Diagrams for Logging/Telemetry Package

This directory contains UML diagrams that document the design and architecture of the logging/telemetry package. These diagrams provide visual representations of the structure, relationships, and interactions between the components of the system.

## Diagrams

### Class Diagram (`class_diagram.md`)

The class diagram illustrates the static structure of the system, showing:
- The core interfaces and classes (`Logger`, `Transaction`, `Driver`)
- The concrete implementations (`logger`, `transaction`, `ConsoleDriver`, etc.)
- The relationships between these components (inheritance, composition, etc.)
- The attributes and methods of each class

This diagram is useful for understanding the overall architecture of the package and how the different components relate to each other.

### Sequence Diagrams

#### Basic Logging Process (`sequence_diagram.md`)

This sequence diagram shows the interactions between components during:
1. Basic logging operations (e.g., `logger.Info()`)
2. Transaction-based logging (e.g., `tx.Info()`)
3. Logger closure (e.g., `logger.Close()`)

It illustrates how log entries flow through the system from the client code to the various output destinations.

#### Configuration Process (`configuration_sequence_diagram.md`)

This sequence diagram focuses on the configuration aspects of the system:
1. Loading configuration from a file
2. Creating drivers based on configuration
3. Constructing a logger with the configured drivers
4. Using the configured logger

### Component Diagram (`component_diagram.md`)

The component diagram provides a high-level view of the system, showing:
- Major components and their grouping into packages
- Dependencies between components
- External interfaces and interactions

This diagram helps understand the system from a deployment and module perspective.

## How to View

These diagrams are written in PlantUML, which is a text-based format for creating UML diagrams. To view them:

1. Use an online PlantUML renderer like [PlantUML Web Server](https://www.plantuml.com/plantuml/uml/)
2. Install a PlantUML plugin for your IDE (e.g., for VS Code, IntelliJ, etc.)
3. Use the PlantUML command-line tool

## Diagram Relationships

- The **Class Diagram** shows the static structure of the code
- The **Sequence Diagrams** show the dynamic behavior of the system
- The **Component Diagram** shows the high-level structure and dependencies

Together, these diagrams provide a comprehensive view of the logging package's architecture and design.

- `component_diagram.md`: High-level components and their relationships (architectural view) 
