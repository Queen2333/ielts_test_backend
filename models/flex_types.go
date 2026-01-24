package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexInt is a custom type that can unmarshal both int and string values from JSON
type FlexInt int

// UnmarshalJSON implements the json.Unmarshaler interface
func (f *FlexInt) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as int first
	var intVal int
	if err := json.Unmarshal(data, &intVal); err == nil {
		*f = FlexInt(intVal)
		return nil
	}

	// If that fails, try to unmarshal as string and convert to int
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return fmt.Errorf("cannot convert string to int: %v", err)
		}
		*f = FlexInt(intVal)
		return nil
	}

	return fmt.Errorf("cannot unmarshal to FlexInt")
}

// MarshalJSON implements the json.Marshaler interface
// Always marshal as int
func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

// Int returns the int value
func (f FlexInt) Int() int {
	return int(f)
}
