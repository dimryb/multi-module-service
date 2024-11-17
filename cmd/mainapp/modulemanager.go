package main

import (
	"encoding/json"
	"fmt"
	"multi-module-service/modules/config"
)

// ModuleManager управляет модулями, их включением и выключением
type ModuleManager struct {
	config *config.Config
}

// NewModuleManager создает новый объект ModuleManager
func NewModuleManager(config *config.Config) *ModuleManager {
	return &ModuleManager{
		config: config,
	}
}

// EnableModule включает модуль
func (m *ModuleManager) EnableModule(module string) error {
	return m.setModuleState(module, true)
}

// DisableModule выключает модуль
func (m *ModuleManager) DisableModule(module string) error {
	return m.setModuleState(module, false)
}

// IsModuleEnabled проверяет, включен ли модуль
func (m *ModuleManager) IsModuleEnabled(module string) (bool, error) {
	value, err := m.config.Get("modules", module)
	if err != nil {
		return false, fmt.Errorf("failed to check module '%s': %w", module, err)
	}

	// Проверяем, является ли значение булевым
	if enabled, ok := value.(bool); ok {
		return enabled, nil
	}

	return false, fmt.Errorf("invalid type for module '%s', expected bool", module)
}

// setModuleState обновляет состояние модуля
func (m *ModuleManager) setModuleState(module string, state bool) error {
	// Получаем текущую секцию модулей
	modules := make(map[string]interface{})
	err := m.config.LoadInto("modules", &modules)
	if err != nil {
		return fmt.Errorf("failed to load 'modules' section: %w", err)
	}

	// Обновляем состояние модуля
	modules[module] = state

	// Применяем изменения обратно в конфигурацию
	modulesJSON, err := json.Marshal(modules)
	if err != nil {
		return fmt.Errorf("failed to marshal updated modules: %w", err)
	}

	// Преобразуем JSON обратно и перезагружаем в конфиг
	err = json.Unmarshal(modulesJSON, &m.config)
	if err != nil {
		return fmt.Errorf("failed to apply updated modules to config: %w", err)
	}

	return nil
}
