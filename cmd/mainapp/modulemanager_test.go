package main

import (
	"testing"

	"multi-module-service/modules/config"
)

func TestModuleManager(t *testing.T) {
	mockConfig := config.NewMockConfig(map[string]interface{}{
		"modules": map[string]bool{
			"module1": true,
			"module2": false,
		},
	})

	manager, err := NewModuleManager(mockConfig)
	if err != nil {
		t.Fatalf("failed to create ModuleManager: %v", err)
	}

	// Проверка включения модуля
	err = manager.EnableModule("module2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Проверка отключения модуля
	err = manager.DisableModule("module1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Проверка отключения всех модулей
	manager.DisableAllModules()
	if manager.modules["module1"] != false || manager.modules["module2"] != false {
		t.Errorf("expected all modules to be disabled")
	}

	// Проверка ошибки при попытке работы с несуществующим модулем
	err = manager.DisableModule("nonexistent")
	if err == nil {
		t.Errorf("expected error for nonexistent module, got nil")
	}
}
