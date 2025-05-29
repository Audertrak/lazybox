# lazybox

## usage

Honestly, the intent for this is mostly quickly surfacing and prepping data for use in prompting large language models, but may be useful for other purposes as well.

### targets

When calling `lazybox` user specify a target to determine what data is being processed. The targets are:

- fs: emit a representation of the filesystem given a path
- file: open and read the contents of a file
- api: parse source code and extract an API
- pkg: crawl a directory and emit a representation of its file/folder structure with relevant metadata and the contents of text files included
- db: fetch data from a database via a middleware

### modes

After selecting a target, the user can specify a mode to determine how the output is formatted. The modes are:

- jsonify: Print a json representation
- commify: Print a comma-separated values representation
- mdify: Print a markdown representation
- tabelify: Print a tabular representation (ala nushell)
- prettify: Print a "pretty" cli representation (ala charmbracelet)
- commentify: Print output to a comment block given a language (e.g. bash, python, etc.)

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

### intermediate representation and pipes

When you run a lazybox command, it will generate an intermediate representation of the data based on the target specified. This intermediate representation is then parsed and transformed to produce the final output in the desired format. Lazybox holds the intermediate representation in memory until the process is complete, which means it is possible to perform multiple operations from the same intermediate representation.

If output from lazybox is piped, lazybox should detect this and continue to hold the intermediate representation from the first command in memory until the full chain is complete. This allows creativity in weaving lazybox commands together with other cli tools or even itself, selectively outputting a raw string or build a queue of outputs/modifiers from the intermediate representation from one command to the next without needing to write it to disk or store it in a variable.

There may be times when you want print multiple output formats or combinations of flags. For example, you may want to pipe the output of `lazybox fs` to a file, then pipe that file into `lazybox file` to read the contents of a specific file in a specific format. Or you may want to have the extra metadata from the verbose flag, but minify the output. Rather than having to write the output to a file or variable, you can simply pipe the output of one lazybox command into another lazybox command, and it will continue to use the intermediate representation from the first command, only transforming the data when the final output is requested.
