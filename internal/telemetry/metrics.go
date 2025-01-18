package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitMetrics() {
	// Configurar exportador de métricas
	exporter, err := otlpmetricgrpc.New(context.Background())
	if err != nil {
		log.Fatalf("falha ao criar o exportador de métricas: %v", err)
	}

	// Configurar o MeterProvider com o exportador
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)

	// Registrar o MeterProvider globalmente
	otel.SetMeterProvider(meterProvider)

	log.Println("Métricas configuradas com sucesso.")
}
