package fileutils

import (
	"bufio"
	"fmt"
	"os"
)

// ReadFile читает файл построчно и возвращает слайс строк
func ReadFile(filename string) ([]string, error) {
	// Открытие файла
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Сканирование строк
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}
