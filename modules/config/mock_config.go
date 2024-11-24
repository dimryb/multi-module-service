package config

import "fmt"

type MockConfig struct {
	data map[string]interface{}
}

func NewMockConfig(initialData map[string]interface{}) *MockConfig {
	return &MockConfig{data: initialData}
}

func (m *MockConfig) LoadInto(key string, target interface{}) error {
	value, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key '%s' not found", key)
	}

	switch t := target.(type) {
	case *map[string]bool:
		if converted, ok := value.(map[string]bool); ok {
			*t = converted
			return nil
		}
	}

	return fmt.Errorf("unsupported type for key '%s'", key)
}
