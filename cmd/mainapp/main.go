package main

// Общий список TODO:
// добавить модуль логирования и мониторинга
// логирование в формате сислога
// добавить модуль конфигурации, выбрать формат конфигурационного файла
// конфигурация модулей фитча-тоглы
// конфигурация в yaml файле
// файлы сервиса, запуск и проверка на винде в его линукс окружении
// настройка пайпа гитлаб, прохождение тестов через пайп
// юнит тесты
// место хранения базовых констант проекта: версия, имя сервиса
// модуль контроля и коррекции топиков в MQTT

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"multi-module-service/modules/config"
	"multi-module-service/modules/mqttclient"
	"multi-module-service/modules/weather"
)

type Main struct {
	appConfig  AppConfig
	mqttConfig MqttConfig
}

type AppConfig struct {
	Weather bool `json:"weather"`
}

type MqttConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (*Main) checkFlags() {
	versionFlag := flag.Bool("version", false, "Показать версию программы")
	debugFlag := flag.Bool("debug", false, "Включить дебаг-режим")
	configFileFlag := flag.String("config", "", "Путь к файлу конфигурации")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s версия %s\n", AppName, Version)
		os.Exit(0)
	}

	if *debugFlag {
		fmt.Println("Дебаг-режим включен")
		// TODO: включить дополнительные логи или отладочную информацию
	}

	if *configFileFlag != "" {
		fmt.Printf("Используется файл конфигурации: %s\n", *configFileFlag)
		// TODO: добавить логику для загрузки иной конфигурации
	}
}

func (m *Main) loadAppConfig() (*config.Config, error) {
	cfg, err := config.NewConfig("config.yml", &config.DefaultReader{})
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Загружаем данные из конфигурации в структуру
	if err := cfg.LoadInto("modules", &m.appConfig); err != nil {
		return nil, fmt.Errorf("loading config into structure: %w", err)
	}

	// Загружаем конфигурацию MQTT
	if err := cfg.LoadInto("mqtt", &m.mqttConfig); err != nil {
		return nil, fmt.Errorf("loading config into structure: %w", err)
	}

	// Выводим данные из конфигурации
	fmt.Println("App config:", m)

	return cfg, nil
}

func main() {
	var m Main

	m.checkFlags()

	cfg, err := m.loadAppConfig()
	if err != nil {
		log.Fatalf("Error loading application config: %v", err)
	}

	// Создаем менеджер модулей
	moduleMgr := NewModuleManager(cfg)

	client, err := mqttclient.NewClient(m.mqttConfig.Host, m.mqttConfig.Port)
	if err != nil {
		log.Fatalf("Ошибка подключения к MQTT: %v", err)
	}
	defer client.Disconnect(250)

	weatherModuleName := "weather"
	enabledWeather, err := moduleMgr.IsModuleEnabled(weatherModuleName)
	if err != nil {
		log.Fatalf("Failed to check module state: %v", err)
	}

	fmt.Printf("Is module '%s' enabled: %v\n", weatherModuleName, enabledWeather)
	if enabledWeather {
		weatherService := weather.NewWeatherService(client, "56.4977", "84.9744", 5*time.Second)
		go weatherService.Run()
	}

	select {}
}
