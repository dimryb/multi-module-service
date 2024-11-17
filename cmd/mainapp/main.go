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

type AppConfig struct {
	Weather bool `json:"weather"`
}

func checkFlags() {
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

func loadAppConfig() (*config.Config, error) {
	cfg, err := config.NewConfig("config.yml")
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// Структура, в которую загружаем конфигурацию
	var appConfig AppConfig

	// Загружаем данные из конфигурации в структуру
	if err := cfg.LoadInto("modules", &appConfig); err != nil {
		return nil, fmt.Errorf("loading config into structure: %w", err)
	}

	// Выводим данные из конфигурации
	fmt.Println("App config:", appConfig)

	return cfg, nil
}

func main() {
	checkFlags()

	cfg, err := loadAppConfig()
	if err != nil {
		log.Fatalf("Error loading application config: %v", err)
	}

	// Создаем менеджер модулей
	moduleMgr := NewModuleManager(cfg)

	client, err := mqttclient.NewClient("192.168.0.6", 1884)
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
