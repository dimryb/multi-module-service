package config

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	data map[string]interface{}
}

// NewConfig создает новый объект конфигурации, автоматически конвертируя YAML или JSON файл в структуру
func NewConfig(filePath string) (*Config, error) {
	// Считываем файл
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config

	// Определяем тип файла по его расширению
	if len(filePath) > 4 && filePath[len(filePath)-4:] == ".yml" || filePath[len(filePath)-5:] == ".yaml" {
		// Если это YAML файл, парсим его как YAML
		if err := yaml.Unmarshal(content, &config.data); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	} else if len(filePath) > 4 && filePath[len(filePath)-5:] == ".json" {
		// Если это JSON файл, парсим его как JSON
		if err := json.Unmarshal(content, &config.data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported config format, must be YAML or JSON")
	}

	return &config, nil
}

// Get получает значение по секции и ключу
func (c *Config) Get(section string, key string) (interface{}, error) {
	if sec, exists := c.data[section]; exists {
		sectionData, ok := sec.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid section type")
		}

		if value, exists := sectionData[key]; exists {
			return value, nil
		}
		return nil, fmt.Errorf("key '%s' not found in section '%s'", key, section)
	}

	return nil, fmt.Errorf("section '%s' not found", section)
}

// LoadInto заполняет переданную структуру значениями из конфигурации
func (c *Config) LoadInto(section string, target interface{}) error {
	// Получаем данные секции из конфигурации
	if sec, exists := c.data[section]; exists {
		sectionData, ok := sec.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid section type")
		}

		// Преобразуем данные секции в JSON
		sectionJSON, err := json.Marshal(sectionData)
		if err != nil {
			return fmt.Errorf("failed to marshal section data: %w", err)
		}

		// Десериализуем JSON в структуру, переданную пользователем
		if err := json.Unmarshal(sectionJSON, target); err != nil {
			return fmt.Errorf("failed to unmarshal section into target structure: %w", err)
		}

		return nil
	}

	return fmt.Errorf("section '%s' not found", section)
}

// GetAll возвращает все данные конфигурации в виде JSON
func (c *Config) GetAll() string {
	jsonData, err := json.Marshal(c.data)
	if err != nil {
		return fmt.Sprintf("Error marshalling config data: %v", err)
	}
	return string(jsonData)
}
