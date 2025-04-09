# Rendering PlantUML Diagrams

The UML diagrams in this directory are written in PlantUML format. To view them as actual diagrams, you can use one of the following methods:

## Option 1: Online PlantUML Server

1. Copy the entire content of any diagram file (e.g., `class_diagram.md`)
2. Go to [PlantUML Web Server](https://www.plantuml.com/plantuml/uml/)
3. Paste the content into the text area
4. The diagram will be rendered automatically

## Option 2: Using Visual Studio Code

1. Install the [PlantUML extension](https://marketplace.visualstudio.com/items?itemName=jebbs.plantuml) for VS Code
2. Open any diagram file
3. Right-click in the editor and select "Preview Current Diagram" or use the keyboard shortcut

## Option 3: Command-line

1. Install PlantUML command-line tool (requires Java):
   ```
   brew install plantuml  # macOS with Homebrew
   apt-get install plantuml  # Ubuntu/Debian
   ```

2. Generate diagrams:
   ```
   plantuml class_diagram.md  # Generates class_diagram.png
   ```

## Troubleshooting Syntax Errors

If you encounter a syntax error such as "Assumed diagram type: class" or "Assumed diagram type: sequence", try these fixes:

1. **Check the file format**:
   - Ensure the file starts with `@startuml` and ends with `@enduml`
   - Remove any extra Markdown formatting (like code fences ```)
   - There should be no text after the `@enduml` tag

2. **Check for syntax issues**:
   - Remove the angle brackets (>) from relationship declarations (e.g., use `..>` instead of `..> >`)
   - Add a title to the diagram with `title Your Title Here` after the `@startuml` tag
   - Check for balanced quotes and parentheses

3. **Simplify if needed**:
   - If a specific part causes issues, try commenting it out with `'` to isolate the problem
   - Split large diagrams into smaller ones
   - Remove complex formatting and add it back gradually

4. **Alternative file formats**:
   - Try saving diagrams with `.puml` extension instead of `.md`
   - Use plain text files without Markdown formatting

## Notes on Diagram Files

Our diagrams follow this naming convention:

- `class_diagram.md`: Static structure of the logging package (classes, interfaces, relationships)
- `sequence_diagram.md`: Sequence of operations for logging (runtime interactions)
- `configuration_sequence_diagram.md`: Sequence for configuration (how configuration works)
- `component_diagram.md`: High-level components and their relationships (architectural view) 
