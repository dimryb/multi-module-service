### README.md

# Конфигурационный модуль

Модуль предоставляет удобный способ загрузки и работы с конфигурацией из файлов **YAML** и **JSON**. Поддерживает извлечение значений по секциям и ключам, а также загрузку данных в структуры Go.

---

## Возможности

- Поддержка форматов **YAML** и **JSON**.
- Извлечение значений по секциям и ключам.
- Загрузка секций в структуры Go.
- Легкая тестируемость через интерфейс чтения файлов.

---

## Установка

Добавьте модуль в проект:

```bash
go get <module_path>/config
```

---

## Примеры использования

### Загрузка конфигурации

```go
reader := &config.DefaultReader{}
cfg, err := config.NewConfig("config.yml", reader)
if err != nil {
	log.Fatalf("Ошибка загрузки конфигурации: %v", err)
}
fmt.Println(cfg.GetAll())
```

---

### Получение значения по ключу

```go
value, err := cfg.Get("section1", "key1")
if err != nil {
	log.Fatalf("Ошибка: %v", err)
}
fmt.Println("Значение:", value)
```

---

### Загрузка секции в структуру

```go
type Section1Config struct {
	Key1 string
	Key2 int
	Key3 bool
}

var section Section1Config
if err := cfg.LoadInto("section1", &section); err != nil {
	log.Fatalf("Ошибка загрузки секции: %v", err)
}
fmt.Printf("Секция: %+v\n", section)
```

---

### Пример конфигурационного файла

#### YAML (`config.yml`)

```yaml
section1:
  key1: "value1"
  key2: 42
  key3: true
section2:
  nestedKey: "nestedValue"
```

#### JSON (`config.json`)

```json
{
  "section1": {
    "key1": "value1",
    "key2": 42,
    "key3": true
  },
  "section2": {
    "nestedKey": "nestedValue"
  }
}
```

---

### Тестирование

Для тестов можно использовать `MockReader`, чтобы эмулировать чтение файлов. Файлы для тестов рекомендуется хранить в директории `testdata`.

Запуск тестов:

```bash
go test ./...
```