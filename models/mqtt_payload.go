package models

type MQTTTelemetryPayload struct {
	AirTemp   float32 `json:"air_temp"`
	WaterTemp float32 `json:"water_temp"`
	Humidity  float32 `json:"humidity"`
}

type MQTTCommandPayload struct {
	Command string `json:"command"`
}
