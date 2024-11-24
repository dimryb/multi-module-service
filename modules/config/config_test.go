package config_test

import (
	"errors"
	"multi-module-service/modules/config"
	"testing"
)

type MockReader struct {
	content []byte
	err     error
}

func (r *MockReader) ReadFile(filePath string) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.content, nil
}

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

func TestNewConfig_SuccessYAML(t *testing.T) {
	mockReader := &MockReader{
		content: []byte(`
modules:
  moduleA: true
  moduleB: false
`),
	}
	cfg, err := config.NewConfig("config.yaml", mockReader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	value, err := cfg.Get("modules", "moduleA")
	if err != nil || value != true {
		t.Fatalf("expected moduleA to be true, got %v, error: %v", value, err)
	}
}

func TestNewConfig_FailureUnsupportedFormat(t *testing.T) {
	mockReader := &MockReader{content: []byte(`{}`)}
	_, err := config.NewConfig("config.txt", mockReader)
	if err == nil {
		t.Fatal("expected error for unsupported format, got none")
	}
}

func TestGet_KeyNotFound(t *testing.T) {
	mockReader := &MockReader{
		content: []byte(`
modules:
  moduleA: true
`),
	}
	cfg, _ := config.NewConfig("config.yaml", mockReader)
	_, err := cfg.Get("modules", "moduleB")
	if err == nil {
		t.Fatal("expected error for missing key, got none")
	}
}

func TestLoadInto_Success(t *testing.T) {
	mockReader := &MockReader{
		content: []byte(`
modules:
  moduleA: true
  moduleB: false
`),
	}
	cfg, _ := config.NewConfig("config.yaml", mockReader)
	var modules map[string]bool
	err := cfg.LoadInto("modules", &modules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !modules["moduleA"] || modules["moduleB"] {
		t.Fatalf("unexpected values in modules: %v", modules)
	}
}

func TestNewConfig_FileReadError(t *testing.T) {
	mockReader := &MockReader{err: errors.New("file not found")}
	_, err := config.NewConfig("config.yaml", mockReader)
	if err == nil {
		t.Fatal("expected error for file read, got none")
	}
}
