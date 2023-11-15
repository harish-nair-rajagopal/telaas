# Helm Chart for Open Telemetry Collector which acts as the global pipeline

## Overview
We are using our custom build collector created as part of [opamp enablement](../opamp/README.md)

The global pipeline acts as the receiver for all agents deployed by different teams. This pipeline will be scaled up based on needs. It processes the incoming data and based on the header `routing-key`, the data is forwarded to the correct exporter.

## Running otel pipeline
From helm\global-pipeline, run
    `helm install global-pipeline open-telemetry/opentelemetry-collector -f global-pipeline.yaml`