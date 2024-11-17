module multi-module-service

go 1.22.2

require (
	multi-module-service/modules/config v0.0.0
	multi-module-service/modules/mqttclient v0.0.0
	multi-module-service/modules/weather v0.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace multi-module-service/modules/config => ./modules/config

replace multi-module-service/modules/mqttclient => ./modules/mqttclient

replace multi-module-service/modules/weather => ./modules/weather
