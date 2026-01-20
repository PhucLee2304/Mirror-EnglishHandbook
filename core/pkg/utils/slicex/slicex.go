package slicex

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func JsonToSlice(jsonData datatypes.JSON) []string {
	var result []string
	if len(jsonData) > 0 {
		_ = json.Unmarshal(jsonData, &result)
	}
	if result == nil {
		return []string{}
	}
	return result
}
