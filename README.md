# Open Telemetry as a Service

Open Telemetry as a Service provides the following

- API/ IaaS bindings for generating open telemetry pipeline based on configuration
- Providing secure endpoints for telemetry data
- Monitoring the telemetry pipeline (owned or externally added) for usage, load, anomalies and security
- Export bindings to single or multiple endpoints for sending data post-processing

## Introduction

Open Telemetry is a popular open framework that is fast gaining popularity and acceptance. In fact, it is the second most active project in CNCF after the core Kubernetes project. Open Telemetry encompasses a set of open protocols, standards, SDKs and tools used for instrumenting, generating, collecting, pre-processing and exporting telemetry data - metrics, traces and logs. It helps in doing away with proprietary telemetry collection and focusses on building open and flexible pipelines that give control of the generated telemetry data to the developers instead of locking it up within silos.

## Architecture


## APIs

### Create Pipeline API
Uses OTEL configuration to create a OTEL collector pipeline which is a pod deployment for this team. Validates configuration for supported receivers, processors and exporters.

Create routing key and add to OTEL Global Pipeline for routing to the right pipeline

reqeust -  OTEL configuration [receiver, processor, exporter]
response - routing key, auth token

### Add Pipeline API
Add an existing OTEL collector running elsewhere to the pipeline management. 

reqeust -  OTEL configuration [receiver, processor, exporter]
response - routing key, auth token


### Delete Pipeline API

request - routing key
response - success/failure

### Update Pipeline API

request - routing key, OTEL configuration [receiver, processor, exporter]
response - routing key, auth token

### Get Pipeline API
request - routing key
response - configuration, pipeline status

## Demo

### Steps

[Building and running the Opamp Server](helm/opamp/README.md)

Building and running the global pipeline service

[Building and running the otaas service](helm/otaas/README.md)

[Building and running the OTEL collector with Opamp enabled](helm/test-collector/README.md)


### Recording

