module multi-module-service/modules/weather

go 1.22.2

replace multi-module-service/modules/mqttclient => ../mqttclient

require multi-module-service/modules/mqttclient v0.0.0-00010101000000-000000000000

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)
