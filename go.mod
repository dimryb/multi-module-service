module multi-module-service

go 1.22.2

require (
	multi-module-service/modules/mqttclient v0.0.0
	multi-module-service/modules/weather v0.0.0
)

replace multi-module-service/modules/mqttclient => ./modules/mqttclient
replace multi-module-service/modules/weather => ./modules/weather
