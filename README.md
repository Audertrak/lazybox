# lazybox

Personal polymorphic structured data CLI swiss army knife.

## usage

Honestly, the intent for this is mostly quickly surfacing and prepping data for use in prompting large language models, but may be useful for other purposes as well.

### targets

When calling `lazybox`, users specify a target to determine what data is being processed. The targets are:

- fs: emit a representation of the filesystem given a path
- file: open and read the contents of a file
- api: parse source code and extract an API
- pkg: crawl a directory and emit a representation of its file/folder structure with relevant metadata and the contents of text files included
- text: parse a text file and extract its contents, including metadata such as word count, line count, and other relevant information
- code: parse source code and extract relevant information, such as functions, classes, and other code constructs
- func: parse a function and extract its signature, parameters, and other relevant information
- env: parse environment variables and emit a representation of their values
- struct: parse a data structure and emit a representation of its contents, including metadata such as field names, types, and other relevant information
- enum: parse an enumeration and emit a representation of its values, including metadata such as field names, types, and other relevant information
- list: parse a data structure and emit a representation of its 'compile time' contents, including metadata such as field names, types, and other relevant information
- db: fetch data from a database via a middleware
- fetch: display system information, similar to fastfetch/neofetch

### modes

After selecting a target, the user can specify a mode to determine how the output is formatted. The modes are:

- jsonify: Print a json representation
- commafy: Print a comma-separated values representation
- mdify: Print a markdown representation
- tabelify: Print a tabular representation (ala nushell)
- prettify: Print a "pretty" cli representation (ala charmbracelet)
- commentify: Print output to a comment block given a language (e.g. bash, python, etc.)
- httpify: Print output as an HTTP response, with appropriate headers and formatting
- flowify: Print output as a flowchart or diagram
- graphify: Print output as a graph or chart
- pdfify: Print output as a PDF document
- astify: Print a structured representation of the data as an abstract syntax tree (AST)
- xmlify: Print output as an XML representation
- structify: Print output to a code struct
- enumify: Print output as an enumeration
- funcify: Print output as a function
- boolify: Print output as boolean logic diagram

> [!NOTE]
> Modes are content sensitive. Some modes may require additional parameters to be specified. For example, the `commentify` mode requires a language to be specified, while the `httpify` mode requires a http verb to be specified. If a mode requires additional parameters, lazybox will prompt the user for those parameters before proceeding with the output.

> [!WARNING]
> Here be dragons. I will attempt to develop lazybox such that it will always _try_ to output something, but there may be times when the output is not what you expect or want. This is especially true for modes that are inherently destructive, or for conversions between typically unrelated or dissimilar data.

### flags

Optionally, the user can specify flags to modify the output further. The flags are:

- all (-a): print all representations of the data, including all available metadata and results
- incremental (-i): print the output incrementally as it is processed, rather than waiting for the entire process to complete; useful for printing multiple representations without needing to print all
- ir (-I): print the intermediate representation of the data, which is a raw, unprocessed version of the data that lazybox holds in memory; useful for debugging or further processing
- less (-l): compact, minimal output, with selective exclusions of metadata or results
- min (-m): remove all whitespace and convert to a single string value
- silent (-s): create an intermediate representation of the data, but do not print it to stdout; useful for piping the output to another command or for debugging
- tokenize (-t): remove articles or other prose grammar and use simple key:value pairs to simplify and shorten output while attempting to preserve meaning.
- verbose (-v): verbose output. Includes additional metadata and results that may not be included in the default output, such as file sizes, line counts, or other relevant information.

___

## Core Architecture: The Universal Data Morphing Framework (UDMF)

Lazybox's power and flexibility stem from its adoption of the Universal Data Morphing Framework (UDMF). At the heart of UDMF is the **Generalized Labeled Property Graph (GLPG)**, which serves as `lazybox`'s primary Intermediate Representation (IR). This approach allows for robust and flexible data manipulation.

### What is a Generalized Labeled Property Graph (GLPG)?

A GLPG is a data structure that represents information as a network of interconnected entities. It consists of:

- **Nodes (or Vertices):** Represent individual objects, items, or concepts (e.g., a file, a function, a database table). Each node usually has a unique identifier and can have labels to categorize it (e.g., "FileInfo", "FunctionInfo").
- **Edges (or Relationships):** Directed connections between two nodes, signifying a relationship (e.g., a "contains" edge from a directory node to a file node, a "calls" edge from one function node to another). Edges have labels to describe the nature of the relationship.
- **Properties:** Key-value pairs that store detailed attributes for both nodes and edges (e.g., a "FileInfo" node might have properties like `name: "report.txt"`, `size: 1024`; an edge might have properties like `type: "dependency"`).

This structure is highly flexible and can model a wide variety of data from different sources in a consistent way.

### The `lazybox` Data Flow with GLPG-IR

The adoption of GLPG as the primary IR refines the data flow within `lazybox`:

1.  **Target Execution:** The user runs `lazybox <target> [args...]`. The specified target's Go package (e.g., `internal/fs`, `internal/code`) executes its logic.
    - _Output:_ This stage produces native Go IR structs, as defined in `internal/ir/ir.go` (e.g., `*ir.FileInfo`, `*ir.CodeInfo`). These structs are tailored to the specific data source.
2.  **GLPG Ingestion:** The native Go IR struct generated by the target is then converted into a GLPG instance by a dedicated "ingestor" component.
    - _Output:_ This in-memory GLPG becomes the **primary Intermediate Representation** that `lazybox` works with for subsequent operations.
3.  **Flag Application (on GLPG - Pre-computation/Filtering):** Many flags (e.g., `--tokenize`, `--less`, `--verbose`) can now operate directly on the GLPG. This might involve transforming the graph (e.g., simplifying its structure, removing or adding nodes/properties) or annotating parts of it for later processing.
    - The `--ir` flag, for instance, will serialize this GLPG directly.
4.  **Mode Serialization:** The selected output mode's handler (e.g., `jsonify`, `mdify`, `prettify`) takes the (potentially transformed) GLPG.
    - It traverses the GLPG and generates the final string output in the desired format. This means modes are essentially GLPG-to-string serializers.
5.  **Final Flag Application (on String Output - Post-computation):** Some flags (e.g., `--min`) might still operate on the final string output generated by the mode serializer to perform simple textual manipulations.

This can be visualized as:

> Target Execution (Source Data -> Native Go IR) -> **GLPG Ingestion** (Native Go IR -> **GLPG**) -> Flag Application (on GLPG) -> Mode Serialization (GLPG -> Formatted String) -> Final Flag Application (on String) -> Output

### Advantages of the GLPG-IR Approach

Using GLPG as the central IR offers significant benefits:

- **Decoupling:** Targets focus on parsing source data into native Go IRs and then ingesting these into a GLPG. Output modes focus on serializing a GLPG into various formats. Targets and modes don't need direct knowledge of each other's specific details, only the common GLPG structure.
- **Universal Transformations:** Complex operations like filtering, enrichment, or structural changes can be defined as transformations on the GLPG, applicable regardless of the original data source.
- **Enhanced Extensibility:**
  - Adding a **new target** involves implementing its data extraction logic (to its native Go IR) and an ingestor to convert that native IR to GLPG.
  - Adding a **new mode** involves implementing a serializer that takes a GLPG and renders it in the new format.
- **Powerful Piping and Chaining:** The GLPG generated by one `lazybox` command can potentially be held in memory and directly consumed or further transformed by subsequent `lazybox` operations or other tools designed to work with a standardized graph representation.
- **Enhanced `--ir` Flag:** The `--ir` flag becomes even more powerful. Instead of just outputting the initial Go struct, it can now output the GLPG itself (e.g., in a standard graph format like JSON Graph Format, GraphML, or a custom verbose representation). This provides a much richer, standardized view of the data `lazybox` is processing and can be consumed by other graph analysis tools.
- **Foundation for Advanced Features:** Sophisticated modes like `flowify` (generating flowcharts) or `graphify` (visualizing data relationships) become more feasible as they can directly operate on the rich structural information within the GLPG.
- **Clear Separation of Concerns:**
  - **Targets:** Source Data -> Native Go IR -> GLPG Ingestor.
  - **Core Logic:** GLPG manipulation and transformation based on flags.
  - **Modes:** GLPG Serializer -> Formatted Output.

### Implications for `lazybox` Usage

- **The Centrality of IR:** The `--ir` flag is crucial for understanding the data `lazybox` is working with, now representing the more universal GLPG.
- **Transformative Power:** Many flags now act as powerful pre-processing steps on this graph-based IR before any textual output is generated.
- **Data Integrity:** While transformations from the source to native Go IR and then to GLPG aim to be faithful, the process of serializing the GLPG to a specific output mode can be additive, selective, subtractive, or transformative. Understanding the GLPG structure helps predict the output.

### Dialects

Structured data could be represented any number of ways. Given that I am developing lazybox primarily for my own use, I have chosen a specific, opinionated output format. That said, the goal is for the project structure to be modular enough to easily extend lazybox to support different preferred output formats or dialects.
This means that while lazybox has a default output format, it is not limited to that format. If somebody finds this project and disagrees with the default format, it should be possible to create their own modules to handle different output formats or dialects, allowing for a wide range of possibilities in how data is represented and processed.

### extensions

Lazybox is designed to be highly modular, and therefore easily extensible, allowing users to add their own targets, modes, and flags. This can be done by creating a new module in the `internal` directory and implementing the necessary functions to handle the target, mode, or flag. The new module should follow the naming conventions and structure of the existing modules to ensure compatibility with lazybox.

#### language support

Certain targets and modes will require awareness of syntax for specific programming languages. I plan on adding support for the languages I use most frequently, which include: C, Go, Lua, Python, and SQL.