package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Reader определяет интерфейс для чтения файлов конфигурации
type Reader interface {
	ReadFile(filePath string) ([]byte, error)
}

// DefaultReader используется для чтения реальных файлов
type DefaultReader struct{}

func (r *DefaultReader) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

type Config struct {
	data map[string]interface{}
}

// NewConfig создает новый объект конфигурации, автоматически конвертируя YAML или JSON файл в структуру
func NewConfig(filePath string, reader Reader) (*Config, error) {
	content, err := reader.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := &Config{data: make(map[string]interface{})}
	if err := parseContent(content, filePath, &config.data); err != nil {
		return nil, err
	}
	return config, nil
}

func parseContent(content []byte, filePath string, data *map[string]interface{}) error {
	if len(filePath) > 4 && (filePath[len(filePath)-4:] == ".yml" || filePath[len(filePath)-5:] == ".yaml") {
		if err := yaml.Unmarshal(content, data); err != nil {
			return fmt.Errorf("failed to parse YAML config: %w", err)
		}
	} else if len(filePath) > 4 && filePath[len(filePath)-5:] == ".json" {
		if err := json.Unmarshal(content, data); err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported config format, must be YAML or JSON")
	}
	return nil
}

// Get получает значение по секции и ключу
func (c *Config) Get(section string, key string) (interface{}, error) {
	sec, exists := c.data[section]
	if !exists {
		return nil, fmt.Errorf("section '%s' not found", section)
	}
	sectionData, ok := sec.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid section type")
	}
	value, exists := sectionData[key]
	if !exists {
		return nil, fmt.Errorf("key '%s' not found in section '%s'", key, section)
	}
	return value, nil
}

// LoadInto заполняет переданную структуру значениями из конфигурации
func (c *Config) LoadInto(section string, target interface{}) error {
	sec, exists := c.data[section]
	if !exists {
		return fmt.Errorf("section '%s' not found", section)
	}
	sectionData, ok := sec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid section type")
	}
	sectionJSON, err := json.Marshal(sectionData)
	if err != nil {
		return fmt.Errorf("failed to marshal section data: %w", err)
	}
	if err := json.Unmarshal(sectionJSON, target); err != nil {
		return fmt.Errorf("failed to unmarshal section into target structure: %w", err)
	}
	return nil
}

// GetAll возвращает все данные конфигурации в виде JSON
func (c *Config) GetAll() string {
	jsonData, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Sprintf("Error marshalling config data: %v", err)
	}
	return string(jsonData)
}
