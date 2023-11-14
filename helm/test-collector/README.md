# Helm Chart for Open Telemetry Collector with OpAmp enabled

## Overview
The default otel collector build from contrib repo does not contain the opamp extension. It has to be build separately.

To do this, we need to follow the steps outlined in https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder

Run the following steps from this directory

- update the custom configuration you need for the otel collector at `test-collector-builder.yaml`

- build a docker image using that binary
    `docker build -f Dockerfile.otel-collector --tag harishrajagopal/otel-collector:0.1.0 .`

    This does the following things
    - go install the builder binary
    - build the custom otel collector binary
    - creates the docker image `harishrajagopal/otel-collector:0.1.0`
    Change the docker image name according to what you need

- push that to your docker repo

- run the otel image using the helm chart or directly using docker run

## Running otel pipeline
From helm\test-collector, run
    `helm install test-otel-collector open-telemetry/opentelemetry-collector -f test-collector.yaml`