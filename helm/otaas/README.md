# Helm chart for Open Telemetry as a service

## Overview

This helm chart deploys otaas on the cluster which provides APIs for creating/updating/deleting OTEL pipelines in the cloud. Any managed pipeline deployed via `/v1/otaas/mPipeline` API is expected to receive otel data from the global pipeline  

For deploying, run the following steps from telaas root
- Build the docker image with the latest changes  
    `docker build -f Dockerfile.otaas --build-arg GITHUB_TOKEN=$GITHUB_TOKEN --tag harishrajagopal/otaas:0.1.0 .`  
    Note: GITHUB_TOKEN needs to be exported for accessing private repos.
- Push the new image to docker hub  
    `docker push harishrajagopal/otaas:0.1.0`  
    Note: Permissions are needed to do this
- Deploy using the helm chart
    `helm install otaas ./helm/otaas/`