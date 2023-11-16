# Helm Chart for installing Open Telemetry Operator

## Overview
We are using our the open telemetry operator for creating otel collectors dynamically through CRDs. Follow instruction from [Open telmetry operator helm chart repo](https://github.com/open-telemetry/opentelemetry-helm-charts/blob/main/charts/opentelemetry-operator/README.md) especially for certmanager installation before running this helm chart. 

The operator needs to be installed before we can utilize the [otaas](../otaas/README.md) APIs

The custom yaml file modifies the default otel collector image used so that it picks our custom build image that has opamp extensions enabled.

## Running open telemetry operator
From helm\otel-operator, run
    `helm install --create-namespace --namespace opentelemetry-operator-system opentelemetry-operator open-telemetry/opentelemetry-operator -f custom-open-telemetry-operator.yaml`