# Helm chart for Open Telemetry Amplifier service to run

## Overview

This helm chart deploys OPAMP (Open Telemetry Amplifier) based service that provides OTEL agent fleet management capabilities

For deploying, run the following steps from telaas root
- Build the docker image with the latest changes  
    `docker build -f Dockerfile.opamp --build-arg GITHUB_TOKEN=$GITHUB_TOKEN --tag harishrajagopal/opamp:0.1.0 .`  
    Note: GITHUB_TOKEN needs to be exported for accessing private repos.
- Push the new image to docker hub  
    `docker push harishrajagopal/opamp:0.1.0`  
    Note: Permissions are needed to do this
- Deploy using the helm chart
    `helm install opamp ./helm/opamp/`

Once the pods are up, you can access the opamp UI from `http://localhost:4321`  

If you launch an OTEL collector that has the opamp agent embedded, as outlined in [test-collector](../test-collector/README.md), the agent will now get listed in the UI. You will also be able to view the configuration details of the opamp agent