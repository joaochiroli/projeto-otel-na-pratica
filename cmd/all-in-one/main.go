// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	configFlag := flag.String("config", "", "path to the config file")
	otelConfigFlag := flag.String("otel", "otel.yaml", "path to the config file")	
	flag.Parse()

	closer, err := telemetry.Setup(context.Background(), *otelConfigFlag)
	if err != nil {
		fmt.Printf("failed to setup telemetry: %v\n", err)
	}
	defer closer(context.Background())

	ctx, span := otel.Tracer("all-in-one").Start(context.Background(), "main")
	

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		otelzap.NewCore("all-in-one", otelzap.WithLoggerProvider(global.GetLoggerProvider())),
	)
	logger := zap.New(core)

	logger.Info("Starting all-in-one service")
	span.AddEvent("Starting all-in-one service")
	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Fatal("failed to load the config", zap.Error(err))
	}

	mux := http.NewServeMux()

	// starts the gRPC server
	lis, err := net.Listen("tcp", c.Server.Endpoint.GRPC)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Fatal("failed to listen", zap.Error(err))
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	{
		logger.Info("Starting the user service")
		span.AddEvent("Starting the user service")
		a := app.NewUser(ctx, &c.Users)
		a.RegisterRoutes(mux)
	}

	{
		logger.Info("Starting the plan service")
		span.AddEvent("Starting the plan service")
		a := app.NewPlan(&c.Plans)
		a.RegisterRoutes(mux, grpcServer)
	}

	{
		span.AddEvent("Starting the payment service")
		logger.Info("Starting the payment service")
		a, err := app.NewPayment(&c.Payments)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.Fatal("failed to create payment service", zap.Error(err))
		}
		a.RegisterRoutes(mux)
		defer func() {
			err = a.Shutdown()
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				logger.Fatal("failed to shutdown payment service", zap.Error(err))
			}
		}()
	}

	{
		logger.Info("Starting the subscription service")
		span.AddEvent("Starting the subscription service")
		a := app.NewSubscription(&c.Subscriptions)
		a.RegisterRoutes(mux)
	}

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			logger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()
	
	span.End()

	err = http.ListenAndServe(c.Server.Endpoint.HTTP, mux)
	if err != nil && err != http.ErrServerClosed {
		logger.Error("failed to serve HTTP", zap.Error(err))
	}

	logger.Info("stopping all-in-one service")
}