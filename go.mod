module github.com/harish-nair-rajagopal/telaas

go 1.15

require github.com/gin-gonic/gin v1.8.1

require github.com/prometheus/client_golang v1.9.0

require (
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.31.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.31.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.8.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.8.0
	go.opentelemetry.io/otel/metric v0.31.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/sdk/metric v0.31.0
	go.opentelemetry.io/otel/trace v1.8.0
	k8s.io/api v0.25.0
	k8s.io/apimachinery v0.25.0
	k8s.io/client-go v0.25.0
//	go.opentelemetry.io/otel/propagation v1.8.0
	github.com/stretchr/testify v1.8.4
	github.com/rs/xid v1.2.1
)
