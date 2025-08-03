package utils

import (
	"encoding/json"
	"fmt"
)

func PrintDebug(name string, data any) {
	if jsonData, err := json.MarshalIndent(data, "", "  "); err == nil {
		fmt.Printf("--- DEBUG [%s] ---\n%s\n", name, string(jsonData))
	} else {
		fmt.Printf("--- DEBUG [%s] ERROR ---\n%v\n", name, err)
	}
}
