package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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

	client, err := mqttclient.NewClient("localhost", 1883)
	if err != nil {
        log.Fatalf("Ошибка подключения к MQTT: %v", err)
    }
    defer client.Disconnect(250)

	weatherService := weather.NewWeatherService(client, "56.4977", "84.9744", 5*time.Second)
	go weatherService.Run()

	select {}
}
