# Open Telemetry as a Service

Open Telemetry as a Service provides the following

- API/ IaaS bindings for generating open telemetry pipeline based on configuration
- Providing secure endpoints for telemetry data
- Monitoring the telemetry pipeline (owned or externally added) for usage, load, anomalies and security
- Export bindings to single or multiple endpoints for sending data post-processing

## Introduction

Open Telemetry is a popular open framework that is fast gaining popularity and acceptance. In fact, it is the second most active project in CNCF after the core Kubernetes project. Open Telemetry encompasses a set of open protocols, standards, SDKs and tools used for instrumenting, generating, collecting, pre-processing and exporting telemetry data - metrics, traces and logs. It helps in doing away with proprietary telemetry collection and focusses on building open and flexible pipelines that give control of the generated telemetry data to the developers instead of locking it up within silos.

Open Agent Management Protocol (OpAMP) is a network protocol for remote management of large fleets of data collection agents. The protocol is vendor agnostic and is currently being adopted into Open Telemetry for managing Open Telemetry agents using extensions in Open Telemetry. This is still WIP at early stage and provides an opportunity for HPE to contribute and guide the vision for the future of observability

## Architecture

![alt text](docs/OTaaS.png)

## Demo

[Presentation](https://hpe-my.sharepoint.com/:p:/p/harish-nair_rajagopal/Ea220lrhevdOh8o-0kvO3ssBW569C4CrKr7JRdrvypcz_g?e=rgTJ2D)

### Steps

[Install open telemetry operator](helm/otel-operator/README.md)

[Building and running the Opamp Server](helm/opamp/README.md)  

[Building and running the global pipeline service](helm/global-pipeline/README.md)  

[Building and running the otaas service](helm/otaas/README.md)  

You can now browse to otaas service UI at http://localhost:4321 and use the UI to manage connected otel collectors and create new ones  

[Building and running the OTEL collector with Opamp enabled](helm/test-collector/README.md)  

The collector launched manually with the above step will also connect to Opamp


### Recording

[Demo Recording 1 - Orchestrating creation of OTEL collector through API](https://hpe-my.sharepoint.com/:v:/p/harish-nair_rajagopal/EeGWUy2yQ3xJliIB0TdCF0MBhrqh9zd2dJPZfGZUwzal8w?e=NYR5ls)

[Demo Recording 2 - Managing launched OTEL collectors via OpAMP service](https://hpe-my.sharepoint.com/:v:/p/harish-nair_rajagopal/Ef6akl_tkWdOnposg-taGbQBdrhiZH9T1I7pUR0UCERJhw?e=eBDpNj)

### References

[OpAmp extensions for Open Telemetry collector](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/opampextension)  
[OpAmp for Open Telemetry Operator](https://docs.google.com/document/d/1M8VLNe_sv1MIfu5bUR5OV_vrMBnAI7IJN-7-IAr37JY/edit#heading=h.bwt48qsb77i2)  
[Current state of OPAMP development in Open Telemetry](https://opentelemetry.io/blog/2023/opamp-status/)  
[Golang implementation of OpAMP](https://github.com/open-telemetry/opamp-go)
