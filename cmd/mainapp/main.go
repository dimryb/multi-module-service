package main

// Общий список TODO:
// добавить модуль логирования и мониторинга
// логирование в формате сислога
// добавить модуль конфигурации, выбрать формат конфигурационного файла
// конфигурация модулей фитча-тоглы
// конфигурация в yaml файле
// документация на каждый модуль
// документация на кросплатформенную сборку
// файлы сервиса, запуск и проверка на винде в его линукс окружении
// настройка пайпа гитлаб, прохождение тестов через пайп
// юнит тесты
// поддержке двух копий репозитория
// место хранения базовых констант проекта: версия, имя сервиса
// модуль контроля и коррекции топиков в MQTT

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"multi-module-service/modules/mqttclient"
	"multi-module-service/modules/weather"
)

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

func main() {
	checkFlags()

	client, err := mqttclient.NewClient("192.168.0.6", 1884)
	if err != nil {
		log.Fatalf("Ошибка подключения к MQTT: %v", err)
	}
	defer client.Disconnect(250)

	weatherService := weather.NewWeatherService(client, "56.4977", "84.9744", 5*time.Second)
	go weatherService.Run()

	select {}
}
