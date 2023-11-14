// Test Application for exploring open telemetry
// Provides metrics, traces ad can generate logs

package main

import (
	"context"
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	//	"go.opentelemetry.io/otel/propagation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

type hw struct {
	greeting string `json:"greeting"`
}

// var hello hw

// Meter can be a global/package variable.
var Meter = global.MeterProvider().Meter("test")

const (
	instrumentationName    = "github.com/cluster-app-test"
	instrumentationVersion = "v0.1.0"
)

var (
	tracer = otel.GetTracerProvider().Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(instrumentationVersion),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
)

// var tracer = otel.Tracer("gin-server")

var dataMap map[string]interface{}

func main() {
	var (
		listen = flag.String("listen", "0.0.0.0:8080", "address of to listen on")
	)
	flag.Parse()
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Create OTLP exporter
	client := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint("test-otel-grpc-pipeline-opentelemetry-collector:3000"),
	)

	exp, err := otlpmetric.New(ctx, client)
	if err != nil {
		log.Fatalf("failed to create the collector exporter: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := exp.Shutdown(ctx); err != nil {
			otel.Handle(err)
		}
	}()
	pusher := controller.New(
		processor.NewFactory(simple.NewWithInexpensiveDistribution(),
			exp,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(2*time.Second),
	)
	global.SetMeterProvider(pusher)
	if err := pusher.Start(ctx); err != nil {
		log.Fatalf("could not start metric controller: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		// pushes any last exports to the receiver
		if err := pusher.Stop(ctx); err != nil {
			otel.Handle(err)
		}
	}()

	gauge, _ := Meter.AsyncFloat64().Gauge(
		"test.cluster.gauge_observer",
		instrument.WithUnit("1"),
		instrument.WithDescription("Test Gauge"),
	)

	if err := Meter.RegisterCallback(
		[]instrument.Asynchronous{
			gauge,
		},
		func(ctx context.Context) {
			gauge.Observe(ctx, rand.Float64())
		},
	); err != nil {
		panic(err)
	}
	if dataMap == nil {
		dataMap = make(map[string]interface{})
	}

	// create new router
	router := gin.Default()
	router.Use(otelgin.Middleware("test-otel-app"))
	router.HandleMethodNotAllowed = true
	// hello = hw{greeting: "Hello World!"}
	// Hello World API
	router.GET("/hello", helloWorld)
	// Status Check for livez readyz endpoint
	router.GET("/status", statusCheck)
	// Metrics endpoint from prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(*listen)
}

// Hello World function always retun status OK and text
func helloWorld(c *gin.Context) {
	hello := hw{greeting: "Hello World!"}
	c.JSON(http.StatusOK, hello)
}

func statusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("otlptrace-example"),
		semconv.ServiceVersionKey.String("0.0.1"),
	)
}

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Printf("error creating OTLP trace exporter")
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
