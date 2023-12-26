package util

import (
	"errors"
	"fmt"
)

func IntValue(value interface{}) (int, error) {
	if intVal, ok := value.(int); ok {
		return intVal, nil
	}
	if floatVal, ok := value.(float32); ok {
		return int(floatVal), nil
	}
	if floatVal, ok := value.(float64); ok {
		return int(floatVal), nil
	}
	return 0, errors.New(fmt.Sprintf("failed to cast %s", value))
}
