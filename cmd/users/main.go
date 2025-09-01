// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"net/http"
	"log"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	texp, err := otlptracehttp.New(context.Background(), otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint("localhost:4318"))
	if err != nil {
		log.Fatalf(err)
	}

	tp := trace.NewTracerProvider(trace.WithResource(resource.Default()), trace.WithBatcher(texp))

	otel.SetTracerProvider(tp)

	// producao
	tr := otel.Tracer("github.com/dosedetelemetria/projeto-otel-na-pratica/users")
	ctx, span := tr.Start(context.Background(), "user.main")
	defer span.End()

	c, _ := config.LoadConfig(*configFlag)

	a := app.NewUser(ctx, &c.Users)
	a.RegisterRoutes(http.DefaultServeMux)
	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
}
