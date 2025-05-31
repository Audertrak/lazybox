package output

import (
	"encoding/csv"
	"fmt"
	"lazybox/internal/glpg"
	"os"
	"sort"
)

// PrintGLPGAsCSV outputs the nodes of a GLPG to CSV format.
// It first prints a header row with "ID", "Label", and all unique property keys found in the nodes.
// Then, it prints one row per node, containing its ID, Label, and corresponding property values.
// If a node does not have a specific property, an empty string is written for that cell.
// Edge information is not currently printed by this function.
func PrintGLPGAsCSV(graph *glpg.GLPG, flags map[string]bool) error {
	if graph == nil || len(graph.Nodes) == 0 {
		// TODO: Use a themed message if we decide to print status messages to stderr
		// For now, CSV output itself is just data to stdout.
		// Consider if "no data" should be an empty CSV or an error/message.
		// For now, returning nil as no error occurred, but no data was written.
		// Alternatively, print a CSV with only headers and no data rows.
		// Let's print headers even if there are no nodes, for consistency.
	}

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Collect all unique property keys from nodes for the header
	propertyKeysMap := make(map[string]struct{})
	for _, node := range graph.Nodes {
		for key := range node.Properties {
			propertyKeysMap[key] = struct{}{}
		}
	}

	// Sort property keys for consistent column order
	sortedPropertyKeys := make([]string, 0, len(propertyKeysMap))
	for key := range propertyKeysMap {
		sortedPropertyKeys = append(sortedPropertyKeys, key)
	}
	sort.Strings(sortedPropertyKeys)

	// Prepare header row
	header := []string{"ID", "Label"}
	header = append(header, sortedPropertyKeys...)

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	if graph == nil || len(graph.Nodes) == 0 {
		// If no nodes, we've written the header, so flush and return.
		writer.Flush()
		return writer.Error() // Check for any error during flush
	}

	// Sort node IDs for consistent row order (optional, but good for diffs/testing)
	nodeIDs := make([]string, 0, len(graph.Nodes))
	for id := range graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	// Write node data rows
	for _, nodeID := range nodeIDs {
		node := graph.Nodes[nodeID]
		row := make([]string, len(header))
		row[0] = node.ID
		// Use the first label if available, or an empty string
		if len(node.Labels) > 0 {
			row[1] = node.Labels[0]
		} else {
			row[1] = ""
		}

		for i, key := range sortedPropertyKeys {
			if val, ok := node.Properties[key]; ok {
				row[i+2] = fmt.Sprintf("%v", val) // +2 because of ID and Label columns
			} else {
				row[i+2] = "" // Empty string for missing properties
			}
		}
		if err := writer.Write(row); err != nil {
			// It might be better to collect errors and continue, or stop on first error.
			// For now, stop on first error.
			return fmt.Errorf("failed to write CSV row for node %s: %w", node.ID, err)
		}
	}

	writer.Flush()
	return writer.Error() // Return any error encountered during writing/flushing
}
