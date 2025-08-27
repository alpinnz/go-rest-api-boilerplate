package helper

import (
	"encoding/json"
	"fmt"
)

// ConvertToSlice mengonversi interface{} menjadi slice bertipe []T
func ConvertToSlice[T any](input interface{}) ([]T, error) {
	// Marshal dulu jadi []byte
	b, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", err)
	}

	// Unmarshal ke slice generic
	var result []T
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return result, nil
}
