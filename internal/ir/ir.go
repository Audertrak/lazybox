package ir

import (
	"time"
)

type FileType string

const (
	FileTypeFile      FileType = "file"
	FileTypeDirectory FileType = "directory"
	FileTypeSymlink   FileType = "symlink"
	FileTypeOther     FileType = "other"
)

type FileInfo struct {
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	AbsolutePath  string            `json:"absolute_path"`
	Type          FileType          `json:"type"`
	Size          int64             `json:"size_bytes"`
	Mode          string            `json:"mode"`
	Owner         string            `json:"owner,omitempty"`
	Group         string            `json:"group,omitempty"`
	ModTime       time.Time         `json:"last_modified"`
	CreateTime    *time.Time        `json:"created,omitempty"`
	SymlinkTarget string            `json:"symlink_target,omitempty"`
	IsGitRepo     bool              `json:"is_git_repo,omitempty"`
	GitRemotes    map[string]string `json:"git_remotes,omitempty"`
	Extension     string            `json:"extension,omitempty"`
	Contents      []*FileInfo       `json:"contents,omitempty"`
	Error         string            `json:"error,omitempty"`
	Content       string            `json:"content,omitempty"`
	// Add optional fields for text analysis
	LineCount int `json:"line_count,omitempty"`
	WordCount int `json:"word_count,omitempty"` // For text files
}

// TextInfo extends FileInfo for more detailed text analysis
type TextInfo struct {
	FileInfo
	CharCount   int                `json:"char_count,omitempty"`
	Encoding    string             `json:"encoding,omitempty"`
	Language    string             `json:"language,omitempty"` // Detected natural language (e.g., "en", "fr")
	Keywords    []KeywordFrequency `json:"keywords,omitempty"` // Top N keywords and their frequencies
	Readability *ReadabilityScores `json:"readability,omitempty"`
	Sentiment   *SentimentScores   `json:"sentiment,omitempty"`
}

type KeywordFrequency struct {
	Keyword   string `json:"keyword"`
	Frequency int    `json:"frequency"`
}

type ReadabilityScores struct {
	FleschKincaidGradeLevel float64 `json:"flesch_kincaid_grade_level,omitempty"`
	// Other scores can be added
}

type SentimentScores struct {
	Polarity     float64 `json:"polarity,omitempty"`     // e.g., -1 (negative) to 1 (positive)
	Subjectivity float64 `json:"subjectivity,omitempty"` // e.g., 0 (objective) to 1 (subjective)
}

// CodePosition represents a line and column in a source file
type CodePosition struct {
	Line   int `json:"line"`
	Column int `json:"column,omitempty"`
}

// CodeLocation represents a span in a source file
type CodeLocation struct {
	FilePath string       `json:"file_path"`
	Start    CodePosition `json:"start_pos"`
	End      CodePosition `json:"end_pos"`
}

// CodeElement is a base struct for common properties of code constructs
type CodeElement struct {
	Name        string        `json:"name"`
	Location    *CodeLocation `json:"location,omitempty"`
	Docstring   string        `json:"docstring,omitempty"`   // Associated documentation/comments
	Annotations []string      `json:"annotations,omitempty"` // e.g., @Override, # noqa
	Language    string        `json:"language,omitempty"`    // Programming language
}

// ParameterInfo represents a function parameter
type ParameterInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	DefaultValue string `json:"default_value,omitempty"`
	IsVariadic   bool   `json:"is_variadic,omitempty"`
}

// FunctionInfo represents a function or method
type FunctionInfo struct {
	CodeElement
	Signature     string          `json:"signature"`
	Parameters    []ParameterInfo `json:"parameters,omitempty"`
	ReturnTypes   []string        `json:"return_types,omitempty"`
	Body          string          `json:"body,omitempty"`       // Optional: full body or summary
	Visibility    string          `json:"visibility,omitempty"` // e.g., "public", "private", "protected"
	IsAsync       bool            `json:"is_async,omitempty"`
	IsStatic      bool            `json:"is_static,omitempty"`
	IsConstructor bool            `json:"is_constructor,omitempty"`
	Throws        []string        `json:"throws,omitempty"` // Types of exceptions/errors thrown
}

// FieldInfo represents a field in a struct or class, or a global variable
type FieldInfo struct {
	CodeElement
	Type         string   `json:"type"`
	DefaultValue string   `json:"default_value,omitempty"`
	Visibility   string   `json:"visibility,omitempty"`
	IsStatic     bool     `json:"is_static,omitempty"`
	Tags         []string `json:"tags,omitempty"` // e.g., `json:"name,omitempty"`
}

// StructInfo represents a struct, class, interface, or trait
type StructInfo struct {
	CodeElement
	Type         string         `json:"struct_type"` // "struct", "class", "interface", "trait"
	Fields       []FieldInfo    `json:"fields,omitempty"`
	Methods      []FunctionInfo `json:"methods,omitempty"`
	InheritsFrom []string       `json:"inherits_from,omitempty"` // Names of parent classes/interfaces
	Implements   []string       `json:"implements,omitempty"`    // Names of implemented interfaces
}

// EnumValueInfo represents a single value within an enumeration
type EnumValueInfo struct {
	Name    string `json:"name"`
	Value   string `json:"value,omitempty"` // Actual value if different from name or explicitly set
	Comment string `json:"comment,omitempty"`
}

// EnumInfo represents an enumeration
type EnumInfo struct {
	CodeElement
	Values []EnumValueInfo `json:"values,omitempty"`
}

// ConstantInfo represents a named constant
type ConstantInfo struct {
	CodeElement
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ImportInfo represents an import statement in a code file
type ImportInfo struct {
	Path  string `json:"path"` // e.g., "fmt", "github.com/user/repo"
	Alias string `json:"alias,omitempty"`
}

// CommentInfo represents a standalone code comment
type CommentInfo struct {
	Location CodeLocation `json:"location"`
	Text     string       `json:"text"`
	IsBlock  bool         `json:"is_block"`
}

// SyntaxErrorInfo represents a syntax error found during parsing
type SyntaxErrorInfo struct {
	Location CodeLocation `json:"location"`
	Message  string       `json:"message"`
}

// CodeInfo represents a parsed source code file
type CodeInfo struct {
	FileInfo
	Language     string            `json:"language"` // e.g., "Go", "Python"
	Imports      []ImportInfo      `json:"imports,omitempty"`
	Functions    []FunctionInfo    `json:"functions,omitempty"`
	Structs      []StructInfo      `json:"structs,omitempty"`
	Enums        []EnumInfo        `json:"enums,omitempty"`
	Constants    []ConstantInfo    `json:"constants,omitempty"`
	GlobalVars   []FieldInfo       `json:"global_variables,omitempty"`
	Comments     []CommentInfo     `json:"comments,omitempty"`
	SyntaxErrors []SyntaxErrorInfo `json:"syntax_errors,omitempty"`
	PackageName  string            `json:"package_name,omitempty"`
	APIs         []ApiInfo         `json:"apis,omitempty"` // If APIs are extracted directly from code
}

// ApiInfo represents an API, potentially derived from code or specs
type ApiInfo struct {
	CodeElement
	Endpoints []EndpointInfo `json:"endpoints,omitempty"` // For HTTP APIs
	// Could also include RPC definitions, etc.
	Version string `json:"version,omitempty"`
	SpecURL string `json:"spec_url,omitempty"` // Link to OpenAPI/Swagger spec
}

type EndpointInfo struct {
	CodeElement
	Path         string          `json:"path"`                            // e.g., "/users/{id}"
	HTTPMethod   string          `json:"http_method"`                     // e.g., "GET", "POST"
	Parameters   []ParameterInfo `json:"parameters,omitempty"`            // Path, query, header params
	RequestBody  string          `json:"request_body_schema,omitempty"`   // Schema or type
	ResponseBody map[int]string  `json:"response_body_schemas,omitempty"` // Status code to schema/type
	Tags         []string        `json:"tags,omitempty"`
}

// DependencyInfo represents a package dependency
type DependencyInfo struct {
	Name        string `json:"name"`                   // e.g., "github.com/spf13/cobra"
	Version     string `json:"version,omitempty"`      // e.g., "v1.7.0"
	Type        string `json:"type,omitempty"`         // e.g., "direct", "indirect", "dev"
	Scope       string `json:"scope,omitempty"`        // e.g., "runtime", "test"
	ResolvedURL string `json:"resolved_url,omitempty"` // URL from which it was resolved, if applicable
	License     string `json:"license,omitempty"`      // SPDX license identifier
	Error       string `json:"error,omitempty"`        // If there was an error fetching info for this dependency
}

// PackageInfo represents a software package
type PackageInfo struct {
	FileInfo                        // Embed FileInfo for path, name, etc. (if it's a local package being analyzed)
	Name          string            `json:"package_name"` // Official package name (e.g., from go.mod, package.json)
	Version       string            `json:"version,omitempty"`
	Description   string            `json:"description,omitempty"`
	Language      string            `json:"language,omitempty"` // e.g., "Go", "JavaScript"
	Dependencies  []DependencyInfo  `json:"dependencies,omitempty"`
	Files         []*FileInfo       `json:"files,omitempty"`       // List of files in the package (can be shallow or deep)
	ReadmePath    string            `json:"readme_path,omitempty"` // Path to README file
	License       string            `json:"license,omitempty"`     // SPDX license identifier
	RepositoryURL string            `json:"repository_url,omitempty"`
	HomepageURL   string            `json:"homepage_url,omitempty"`
	EntryPoint    string            `json:"entry_point,omitempty"`   // Main executable or entry file
	BuildScripts  map[string]string `json:"build_scripts,omitempty"` // e.g., "build": "go build ."
	Error         string            `json:"error,omitempty"`
}

// EnvironmentVariable represents a single environment variable
type EnvironmentVariable struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Source   string `json:"source,omitempty"`    // e.g., ".env file", "shell", "system"
	IsSecret bool   `json:"is_secret,omitempty"` // Heuristic to guess if it's a secret
}

// EnvironmentInfo contains a list of environment variables
type EnvironmentInfo struct {
	Variables []EnvironmentVariable `json:"variables"`
	Count     int                   `json:"count"`
	Error     string                `json:"error,omitempty"`
}

// ListItem represents an item within a ListInfo
// It's a generic wrapper to allow for mixed types if necessary,
// though often lists will be homogeneous.
type ListItem struct {
	Index int         `json:"index"`
	Value interface{} `json:"value"`          // Can be any type: string, number, bool, or another IR struct
	Type  string      `json:"type,omitempty"` // Detected type of the value
}

// ListInfo represents a list or array, typically from "compile-time" data
// This is for representing generic lists where items might be of mixed types or simple scalars.
type ListInfo struct {
	CodeElement            // If the list is defined in code, it might have a name, location, etc.
	Items       []ListItem `json:"items"`
	ItemCount   int        `json:"item_count"`
	ListType    string     `json:"list_type,omitempty"`   // e.g., "array", "slice", "tuple"
	Homogeneous bool       `json:"homogeneous,omitempty"` // Whether all items are of the same type
	Error       string     `json:"error,omitempty"`
}

// --- Database Related IRs ---

// ColumnInfo describes a database column
type ColumnInfo struct {
	Name            string `json:"name"`
	Type            string `json:"type"` // e.g., "VARCHAR(255)", "INT", "BOOLEAN"
	IsNullable      bool   `json:"is_nullable"`
	DefaultValue    string `json:"default_value,omitempty"`
	IsPrimaryKey    bool   `json:"is_primary_key,omitempty"`
	IsForeignKey    bool   `json:"is_foreign_key,omitempty"`
	References      string `json:"references,omitempty"` // Table(column) it references if FK
	Comment         string `json:"comment,omitempty"`
	Collation       string `json:"collation,omitempty"`
	CharMaxLength   int    `json:"char_max_length,omitempty"`
	OrdinalPosition int    `json:"ordinal_position"`
}

// IndexInfo describes a database index
type IndexInfo struct {
	Name         string   `json:"name"`
	Columns      []string `json:"columns"`
	IsUnique     bool     `json:"is_unique"`
	IsPrimaryKey bool     `json:"is_primary_key,omitempty"` // If this index is the PK
	Type         string   `json:"type,omitempty"`           // e.g., "BTREE", "HASH"
	Definition   string   `json:"definition,omitempty"`     // SQL definition of the index
}

// ForeignKeyConstraintInfo describes a foreign key constraint
type ForeignKeyConstraintInfo struct {
	Name                string   `json:"name"`
	SourceTable         string   `json:"source_table"`
	SourceColumns       []string `json:"source_columns"`
	TargetTable         string   `json:"target_table"`
	TargetColumns       []string `json:"target_columns"`
	OnUpdate            string   `json:"on_update,omitempty"` // e.g., "CASCADE", "SET NULL"
	OnDelete            string   `json:"on_delete,omitempty"`
	MatchOption         string   `json:"match_option,omitempty"` // e.g., "SIMPLE", "FULL"
	IsDeferrable        bool     `json:"is_deferrable,omitempty"`
	IsInitiallyDeferred bool     `json:"is_initially_deferred,omitempty"`
}

// TableInfo describes a database table
type TableInfo struct {
	Name        string                     `json:"name"`
	Schema      string                     `json:"schema"`
	Columns     []ColumnInfo               `json:"columns,omitempty"`
	Indexes     []IndexInfo                `json:"indexes,omitempty"`
	ForeignKeys []ForeignKeyConstraintInfo `json:"foreign_keys,omitempty"`
	Comment     string                     `json:"comment,omitempty"`
	RowCount    int64                      `json:"row_count,omitempty"`  // Approximate row count
	CreateSQL   string                     `json:"create_sql,omitempty"` // CREATE TABLE statement
	Type        string                     `json:"type,omitempty"`       // e.g., "BASE TABLE", "VIEW"
}

// SchemaInfo describes a database schema
type SchemaInfo struct {
	Name          string      `json:"name"`
	Tables        []TableInfo `json:"tables,omitempty"`
	Views         []TableInfo `json:"views,omitempty"`     // Views can also be represented by TableInfo
	Functions     []string    `json:"functions,omitempty"` // Names of functions/procedures
	Owner         string      `json:"owner,omitempty"`
	Comment       string      `json:"comment,omitempty"`
	CollationName string      `json:"collation_name,omitempty"`
	CharSetName   string      `json:"character_set_name,omitempty"`
}

// QueryResultInfo represents the result of a database query
type QueryResultInfo struct {
	Query           string          `json:"query"`
	Columns         []string        `json:"columns,omitempty"` // Column names
	Rows            [][]interface{} `json:"rows,omitempty"`    // Slice of rows, each row is a slice of values
	RowCount        int             `json:"row_count"`
	Error           string          `json:"error,omitempty"`
	ExecutionTimeMs int64           `json:"execution_time_ms,omitempty"`
	AffectedRows    int64           `json:"affected_rows,omitempty"` // For DML statements
	Messages        []string        `json:"messages,omitempty"`      // e.g., notices from the DB
}
