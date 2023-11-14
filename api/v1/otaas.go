package v1

// OTaaSPipeline is to create otel pipeline in k8s cluster.
type OTaaSPipeline struct {
	Name     string `json:"name"`
	Exporter string `json:"exporter"`
}

// OTaaSRes is response for otel pipeline.
type OTaaSRes struct {
	RoutingKey string `json:"routingkey"`
}
