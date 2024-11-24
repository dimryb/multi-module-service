package main

import (
	"fmt"
	"multi-module-service/modules/config"
)

// ModuleManager управляет модулями, их включением и выключением
type ModuleManager struct {
	config  config.ConfigLoader
	modules map[string]bool
}

// NewModuleManager создает новый объект ModuleManager
func NewModuleManager(config config.ConfigLoader) (*ModuleManager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Инициализация модулей
	modules := make(map[string]bool)
	err := config.LoadInto("modules", &modules)
	if err != nil {
		return nil, fmt.Errorf("failed to load modules: %w", err)
	}

	return &ModuleManager{
		config:  config,
		modules: modules,
	}, nil
}

// ReloadModules считывает состояние модулей из конфигурации и обновляет локальное состояние
func (m *ModuleManager) ReloadModules() error {
	modules := make(map[string]bool)
	err := m.config.LoadInto("modules", &modules)
	if err != nil {
		return fmt.Errorf("failed to load 'modules' section: %w", err)
	}
	m.modules = modules
	return nil
}

// EnableModule включает модуль
func (m *ModuleManager) EnableModule(module string) error {
	if _, exists := m.modules[module]; !exists {
		return fmt.Errorf("module '%s' not found", module)
	}
	m.modules[module] = true
	return nil
}

// DisableModule выключает модуль
func (m *ModuleManager) DisableModule(module string) error {
	if _, exists := m.modules[module]; !exists {
		return fmt.Errorf("module '%s' not found", module)
	}
	m.modules[module] = false
	return nil
}

// DisableAllModules выключает все модули
func (m *ModuleManager) DisableAllModules() {
	for module := range m.modules {
		m.modules[module] = false
	}
}

// IsModuleEnabled проверяет, включен ли модуль
func (m *ModuleManager) IsModuleEnabled(module string) (bool, error) {
	state, exists := m.modules[module]
	if !exists {
		return false, fmt.Errorf("module '%s' not found", module)
	}
	return state, nil
}

// GetModules возвращает текущее состояние всех модулей
func (m *ModuleManager) GetModules() map[string]bool {
	return m.modules
}
