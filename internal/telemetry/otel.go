package telemetry

func InitTelemetry() {
	InitTraces()
	InitMetrics()
	InitLogs()
}
