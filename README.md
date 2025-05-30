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

## Concepts

### intermediate representation and pipes

When you run a lazybox command, it will generate an intermediate representation of the data based on the target specified. This intermediate representation is then parsed and transformed to produce the final output in the desired format. Lazybox holds the intermediate representation in memory until the process is complete, which means it is possible to perform multiple operations from the same intermediate representation.

If output from lazybox is piped, lazybox should detect this and continue to hold the intermediate representation from the first command in memory until the full chain is complete. This allows creativity in weaving lazybox commands together with other cli tools or even itself, selectively outputting a raw string or building a queue of outputs/modifiers from the intermediate representation from one command to the next without needing to write it to disk or store it in a variable.

There may be times when you want print multiple output formats or combinations of flags. For example, you may want to pipe the output of `lazybox fs` to a file, then pipe that file into `lazybox file` to read the contents of a specific file in a specific format. Or you may want to have the extra metadata from the verbose flag, but minify the output. Rather than having to write the output to a file or variable, you can simply pipe the output of one lazybox command into another lazybox command, and it will continue to use the intermediate representation from the first command, only transforming the data when the final output is requested.

### structured data

Depending on the mode and flags, then the resulting _output_ may be:

- additive
- selective
- subtractive
- transformative
- destructive

based on the generated intermediate representation. In order to avoid data loss or corruption, it is important to understand how data is input, parsed, and output. The general flow of data through lazybox is as follows:

> target -> IR -> flags -> mode -> flags -> output

The intermediate representation that is generated by lazybox is a function of the specified target. Certain flags act as filters or modify what IR is then interpreted. The mode selection determines how the IR is 'hydrated' into the final output. Finally, some flags take the hydrated modal output and perform filtering or transformation.

Each target has its own IR structure and format. Conversion from IR to any output is a transformative operation: so as long as the IR is retained, there is no loss of information. However, if the IR is not retained, then the output may be destructive, meaning that some information may be lost or altered in the process.

Furthermore, while it may be possible to consume structured data via pipes or as a direct target, some flags or output modes may be inherently destructive and mangle the payload such that it cannot be recovered or easily converted unless the IR is retained.

If you need to retain the original structure of the data, it is recommended to use the `ir` flag to print the intermediate representation, which will retain the original structure of the data.

### Dialects

Structured data could be represented any number of ways. Given that I am developing lazybox primarily for my own use, I have chosen a specific, opinionated output format. That said, the goal is for the project structure to be modular enough to easily extend lazybox to support different preferred output formats or dialects.
This means that while lazybox has a default output format, it is not limited to that format. If somebody finds this project and disagrees with the default format, it should be possible to create their own modules to handle different output formats or dialects, allowing for a wide range of possibilities in how data is represented and processed.

### extensions

Lazybox is designed to be highly modular, and therefore easily extensible, allowing users to add their own targets, modes, and flags. This can be done by creating a new module in the `internal` directory and implementing the necessary functions to handle the target, mode, or flag. The new module should follow the naming conventions and structure of the existing modules to ensure compatibility with lazybox.

#### language support

Certain targets and modes will require awareness of syntax for specific programming languages. I plan on adding support for the languages I use most frequently, which include: C, Go, Lua, Python, and SQL.