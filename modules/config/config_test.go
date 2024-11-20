package config_test

import (
	"multi-module-service/modules/config"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		filePath  string
		expectErr bool
	}{
		{"Valid YAML", "./testdata/config_valid.yaml", false},
		{"Valid JSON", "./testdata/config_valid.json", false},
		{"Invalid YAML", "./testdata/config_invalid.yaml", true},
		{"Invalid JSON", "./testdata/config_invalid.json", true},
		{"Missing Section", "./testdata/config_missing_section.yaml", false},
		{"Missing Key", "./testdata/config_missing_key.yaml", false},
		{"Edge Cases YAML", "./testdata/config_edgecase.yaml", false},
		{"Edge Cases JSON", "./testdata/config_edgecase.json", false},
		{"Large Config YAML", "./testdata/config_large.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := config.NewConfig(tt.filePath, &config.DefaultReader{})
			if (err != nil) != tt.expectErr {
				t.Errorf("NewConfig() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
