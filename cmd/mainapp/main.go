package main

import (
	"flag"
	"fmt"
	"os"

	"multi-module-service/modules/mqttclient"
	"multi-module-service/modules/weather"
)

const (
	Version = "2.0.0"
)

func checkFlags() {
	versionFlag := flag.Bool("version", false, "Показать версию программы")
	debugFlag := flag.Bool("debug", false, "Включить дебаг-режим")
	configFileFlag := flag.String("config", "", "Путь к файлу конфигурации")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-weather-mqtt версия %s\n", Version)
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

	mqttClient := mqttclient.NewClient("localhost", 1883)

	weatherService := weather.NewWeatherService("56.4977", "84.9744", mqttClient, 5)
	weatherService.StartWeatherCycle()

	select {}
}
