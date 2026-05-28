package models

type MQTTTelemetryPayload struct {
	Event           string  `json:"event"`
	AirTemp         float32 `json:"air_temp"`
	WaterTemp       float32 `json:"water_temp"`
	Humidity        float32 `json:"humidity"`
	WaterSufficient bool    `json:"water_sufficient"`
}

type MQTTCommandPayload struct {
	Command string `json:"command"`
}

type MQTTErrorPayload struct {
	Event  string   `json:"event"`
	Fields []string `json:"fields,omitempty"`
}
