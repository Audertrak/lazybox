package output

import (
	"encoding/json"
	"fmt"
	"lazybox/internal/ir"
)

func PrintJSON(info *ir.FileInfo) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}
