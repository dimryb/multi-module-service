package fileutils

import (
	"os"
	"testing"
)

// TestReadFile_Success проверяет успешное чтение файла
func TestReadFile_Success(t *testing.T) {
	// Подготовка: создание временного файла с тестовыми данными
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после теста

	// Запись данных в файл
	content := "line1\nline2\nline3"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Закрытие файла перед чтением
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Тестируем ReadFile
	lines, err := ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Проверка результата
	expectedLines := []string{"line1", "line2", "line3"}
	if len(lines) != len(expectedLines) {
		t.Fatalf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}
	for i, line := range lines {
		if line != expectedLines[i] {
			t.Errorf("Expected line %q, got %q", expectedLines[i], line)
		}
	}
}

// TestReadFile_FileNotFound проверяет чтение несуществующего файла
func TestReadFile_FileNotFound(t *testing.T) {
	_, err := ReadFile("nonexistentfile.txt")
	if err == nil {
		t.Fatal("Expected an error for nonexistent file, got nil")
	}
}

// TestReadFile_EmptyFile проверяет чтение пустого файла
func TestReadFile_EmptyFile(t *testing.T) {
	// Подготовка: создание временного пустого файла
	tempFile, err := os.CreateTemp("", "emptyfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Удалить файл после теста

	// Тестируем ReadFile
	lines, err := ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Проверка результата
	if len(lines) != 0 {
		t.Fatalf("Expected 0 lines, got %d", len(lines))
	}
}
