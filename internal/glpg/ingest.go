package glpg

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"lazybox/internal/ir" // Corrected import path

	"github.com/google/uuid"
)

// ToGLPG converts any supported IR struct into a GLPG.
// This function acts as a dispatcher based on the type of the input IR.
func ToGLPG(data interface{}) (*GLPG, error) {
	g := NewGLPG()
	err := ingestToGLPG(data, g, "", "") // No parent node or edge label for the root
	if err != nil {
		return nil, err
	}
	return g, nil
}

// ingestToGLPG is the core recursive function that converts IR structs to GLPG nodes and edges.
// parentNodeID and edgeLabelToParent are used to link child nodes to their parent.
func ingestToGLPG(data interface{}, g *GLPG, parentNodeID string, edgeLabelToParent string) error {
	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		return nil // Skip nil data
	}

	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle slices by iterating and ingesting each element
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			// Each element in the slice will be connected to the same parentNodeID
			// with the same edgeLabelToParent, but as a new node.
			err := ingestToGLPG(val.Index(i).Interface(), g, parentNodeID, edgeLabelToParent)
			if err != nil {
				// Potentially collect errors and continue, or return immediately
				return fmt.Errorf("error ingesting slice element %d: %w", i, err)
			}
		}
		return nil
	}

	// Expect a struct at this point for creating a node
	if val.Kind() != reflect.Struct {
		// If we have a parent and an edge label, but the data is a primitive,
		// it might be a property of the parent. However, our design is to make nodes from structs.
		// This case might indicate an unexpected data structure or a need to handle direct properties.
		// For now, we skip non-structs that are not part of a slice handled above.
		// If parentNodeID is not empty, this primitive could have been a property of the parent.
		// However, the main loop over struct fields handles properties.
		return nil
	}

	nodeID := generateNodeID(val.Type().Name(), val)
	node := &GLPGNode{
		ID:         nodeID,
		Labels:     []string{val.Type().Name()},
		Properties: make(GLPGProperty),
	}

	// Add properties from struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Handle embedded structs (only if the embedded struct itself is an exported type)
		if field.Anonymous && fieldVal.Kind() == reflect.Struct && fieldVal.Type().PkgPath() == "" {
			// Recursively add properties from exported fields of the embedded struct
			for j := 0; j < fieldVal.NumField(); j++ {
				embeddedField := fieldVal.Type().Field(j)
				embeddedFieldVal := fieldVal.Field(j)
				// Skip unexported fields in embedded struct
				if embeddedField.PkgPath != "" {
					continue
				}
				if embeddedFieldVal.CanInterface() {
					addProperty(node.Properties, embeddedField.Name, embeddedFieldVal.Interface())
				}
			}
			continue // Done with this embedded struct field
		}

		kind := fieldVal.Kind()
		fieldType := fieldVal.Type()

		// Determine if the field's type is an exported struct or pointer to an exported struct
		isExportedStructType := false
		if kind == reflect.Struct && fieldType.PkgPath() == "" {
			isExportedStructType = true
		} else if kind == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct && fieldType.Elem().PkgPath() == "" {
			isExportedStructType = true
		}

		// If field is an exported struct, pointer to an exported struct, or a slice,
		// it might become new node(s) or its elements processed.
		if isExportedStructType || kind == reflect.Slice {
			if fieldVal.CanInterface() {
				err := ingestToGLPG(fieldVal.Interface(), g, nodeID, field.Name)
				if err != nil {
					return fmt.Errorf("error ingesting field %s: %w", field.Name, err)
				}
			}
		} else {
			// Otherwise, add as a property (e.g. basic types, or non-pointer, non-slice structs
			// that are not meant to be separate nodes but rather flattened - addProperty will handle this).
			if fieldVal.CanInterface() {
				addProperty(node.Properties, field.Name, fieldVal.Interface())
			}
		}
	}

	g.AddNode(node)

	// If this node has a parent, create an edge to it
	if parentNodeID != "" && edgeLabelToParent != "" {
		edge := &GLPGEdge{
			ID:         uuid.NewString(), // Simple unique ID for the edge
			SourceID:   parentNodeID,
			TargetID:   nodeID,
			Label:      edgeLabelToParent,
			Properties: make(GLPGProperty), // Edges can also have properties if needed
		}
		g.AddEdge(edge)
	}

	return nil
}

// addProperty adds a value to the properties map, handling common types.
func addProperty(props GLPGProperty, key string, value interface{}) {
	refVal := reflect.ValueOf(value) // Value is from fieldVal.Interface()

	switch refVal.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		props[key] = value
	default:
		// Special handling for time.Time
		if t, ok := value.(time.Time); ok {
			props[key] = t.Format(time.RFC3339)
		} else if refVal.Kind() == reflect.Struct {
			// This handles structs that are direct fields (not pointers/slices that become nodes)
			// and are intended to be flattened as properties (e.g., ir.ReadabilityScores).
			// Ensure the struct type itself is exported.
			if refVal.Type().PkgPath() == "" { // Exported struct type
				structProps := make(map[string]interface{})
				for i := 0; i < refVal.NumField(); i++ {
					structField := refVal.Type().Field(i)
					structFieldVal := refVal.Field(i)

					// Only consider exported fields of this struct
					if structField.PkgPath == "" && structFieldVal.CanInterface() {
						fieldInterface := structFieldVal.Interface()
						// Handle time.Time within these structs specifically
						if t, ok := fieldInterface.(time.Time); ok {
							structProps[structField.Name] = t.Format(time.RFC3339)
						} else {
							// Add only if it's a basic type to prevent deep nesting issues.
							switch structFieldVal.Kind() {
							case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
								reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
								structProps[structField.Name] = fieldInterface
								// Other nested complex types within these flattened structs are skipped.
							}
						}
					}
				}
				if len(structProps) > 0 {
					props[key] = structProps
				}
			}
		} else if refVal.Kind() == reflect.Map {
			// Handle maps (simplistic: string keys, basic value types)
			mapInterface := make(map[string]interface{})
			iter := refVal.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				if k.CanInterface() && v.CanInterface() {
					// Ensure value is a basic type or time.Time for simplicity in properties
					valInterface := v.Interface()
					if t, ok := valInterface.(time.Time); ok {
						mapInterface[fmt.Sprintf("%v", k.Interface())] = t.Format(time.RFC3339)
					} else {
						switch v.Kind() {
						case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
							reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
							mapInterface[fmt.Sprintf("%v", k.Interface())] = valInterface
						}
					}
				}
			}
			if len(mapInterface) > 0 {
				props[key] = mapInterface
			}
		}
		// Other complex types not explicitly handled are skipped as properties.
	}
}

// generateNodeID creates a unique ID for a node.
// This is a simple implementation and might need to be more robust
// to ensure global uniqueness if graphs are merged, etc.
// For ir.FileInfo, using the path is a good candidate for uniqueness.
func generateNodeID(typeName string, val reflect.Value) string {
	// Attempt to find a unique identifier field like 'Path' or 'Name'
	if val.Kind() == reflect.Struct {
		pathField := val.FieldByName("Path")
		if pathField.IsValid() && pathField.Kind() == reflect.String && pathField.String() != "" {
			// Sanitize path to be a valid ID component
			cleanPath := strings.ReplaceAll(filepath.ToSlash(pathField.String()), "/", "_")
			cleanPath = strings.ReplaceAll(cleanPath, ":", "") // Remove colons for Windows paths
			return fmt.Sprintf("%s_%s", typeName, cleanPath)
		}
		nameField := val.FieldByName("Name")
		if nameField.IsValid() && nameField.Kind() == reflect.String && nameField.String() != "" {
			return fmt.Sprintf("%s_%s", typeName, nameField.String())
		}
	}
	// Fallback to UUID if no suitable field is found or not a struct
	return fmt.Sprintf("%s_%s", typeName, uuid.NewString())
}

// Specific Ingestors (can be added for more control if needed, but generic one is powerful)

// FileInfoToGLPG converts an ir.FileInfo struct and its children into GLPG nodes and edges.
// This is an example of a specific ingestor, though the generic one aims to handle this.
func FileInfoToGLPG(fi *ir.FileInfo, g *GLPG, parentNodeID string, edgeLabel string) error {
	if fi == nil {
		return nil
	}

	nodeID := generateNodeID("FileInfo", reflect.ValueOf(*fi)) // Use path for FileInfo ID
	node := &GLPGNode{
		ID:     nodeID,
		Labels: []string{"FileInfo", string(fi.Type)}, // Add FileType as a label
		Properties: GLPGProperty{
			"Name":         fi.Name,
			"Path":         fi.Path,
			"AbsolutePath": fi.AbsolutePath,
			"Type":         string(fi.Type),
			"Size":         fi.Size,
			"Mode":         fi.Mode,
			"ModTime":      fi.ModTime.Format(time.RFC3339),
			"IsDir":        fi.IsDir,     // Corrected: Direct field access
			"IsSymlink":    fi.IsSymlink, // Corrected: Direct field access
		},
	}

	if fi.SymlinkTarget != "" {
		node.Properties["SymlinkTarget"] = fi.SymlinkTarget
	}
	if !fi.CreateTime.IsZero() {
		node.Properties["CreateTime"] = fi.CreateTime.Format(time.RFC3339)
	}
	if fi.Owner != "" {
		node.Properties["Owner"] = fi.Owner
	}
	if fi.Group != "" {
		node.Properties["Group"] = fi.Group
	}
	if fi.Extension != "" {
		node.Properties["Extension"] = fi.Extension
	}
	if fi.Error != "" {
		node.Properties["Error"] = fi.Error
	}

	// Handle content: if present, add as property. Could also be a separate node for large content.
	if fi.Content != nil && *fi.Content != "" {
		// Decide on a strategy: store directly, or hash, or summarize for large content
		// For now, storing directly if not too large, otherwise a placeholder or hash.
		const maxContentLength = 1024 // Example limit
		if len(*fi.Content) > maxContentLength {
			node.Properties["ContentSummary"] = (*fi.Content)[:maxContentLength] + "... (truncated)"
		} else {
			node.Properties["Content"] = *fi.Content
		}
	}

	if fi.TextAnalysis != nil {
		// Flatten TextAnalysis properties with a prefix
		node.Properties["TextAnalysis_LineCount"] = fi.TextAnalysis.LineCount
		node.Properties["TextAnalysis_WordCount"] = fi.TextAnalysis.WordCount
		node.Properties["TextAnalysis_CharCount"] = fi.TextAnalysis.CharCount
		if fi.TextAnalysis.DetectedLanguage != "" {
			node.Properties["TextAnalysis_DetectedLanguage"] = fi.TextAnalysis.DetectedLanguage
		}
		// Add other TextAnalysis fields as needed, checking for nil pointers if applicable
		if fi.TextAnalysis.Readability != nil {
			node.Properties["TextAnalysis_Readability_FleschKincaidGradeLevel"] = fi.TextAnalysis.Readability.FleschKincaidGradeLevel
			node.Properties["TextAnalysis_Readability_GunningFogIndex"] = fi.TextAnalysis.Readability.GunningFogIndex
		}
		if fi.TextAnalysis.Sentiment != nil {
			node.Properties["TextAnalysis_Sentiment_Polarity"] = fi.TextAnalysis.Sentiment.Polarity
			node.Properties["TextAnalysis_Sentiment_Subjectivity"] = fi.TextAnalysis.Sentiment.Subjectivity
		}
	}

	if fi.GitRemoteURL != "" {
		node.Properties["GitRemoteURL"] = fi.GitRemoteURL
	}
	if fi.GitCurrentBranch != "" {
		node.Properties["GitCurrentBranch"] = fi.GitCurrentBranch
	}

	// Add any custom metadata, prefixing keys to avoid collisions
	for k, v := range fi.Metadata {
		propsKey := "Metadata_" + k
		// Ensure the value is suitable for GLPGProperty (e.g., basic types, strings)
		// This might require more sophisticated type handling if metadata contains complex structs
		node.Properties[propsKey] = fmt.Sprintf("%v", v) // Simple conversion to string
	}

	g.AddNode(node)

	if parentNodeID != "" && edgeLabel != "" {
		edge := &GLPGEdge{
			ID:       uuid.NewString(),
			SourceID: parentNodeID,
			TargetID: nodeID,
			Label:    edgeLabel,
		}
		g.AddEdge(edge)
	}

	// Recursively process children
	for _, child := range fi.Children {
		// Children are linked to the current node with an edge type like "CONTAINS" or "CHILD_OF"
		err := FileInfoToGLPG(child, g, nodeID, "CONTAINS")
		if err != nil {
			return fmt.Errorf("failed to ingest child FileInfo for %s: %w", child.Name, err)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
