package app

// EventStore is an interface for storing events
type EventStore interface {
	// SaveURLHashTelemetry saves url hash telemetry data
	SaveURLHashTelemetry(params URLHashTelemetryParams) error

	// ExportURLHashTelemetry exports all url hash telemetry data
	ExportURLHashTelemetry() ([]URLHashTelemetryParams, error)
}
