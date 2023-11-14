package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"log"

	v1 "github.com/harish-nair-rajagopal/telaas/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	var (
		listen = flag.String("listen", "0.0.0.0:8080", "address of to listen on")
	)
	flag.Parse()
	// ctx := context.Background()

	// create new router
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	// GET pipeline

	// GET pipelines

	// POST managed pipeline (create/ update)
	router.POST("/v1/otaas/mPipeline", CreateOTELPipeline)
	// POST unmanaged pipeline (add)

	// DELETE pipeline

	router.Run(*listen)
}

func CreateOTELPipeline(g *gin.Context) {
	ctx := g.Request.Context()

	var otaasPipeLine v1.OTaaSPipeline

	if err := g.BindJSON(&otaasPipeLine); err != nil {
		log.Printf(ctx, "error: %v", err)
		g.JSON(http.StatusBadRequest, err.Error())

		return
	}

	log.Printf(ctx, "creating OtelPipe with details [%+v]", otaasPipeLine)

	//Construct Custom Resource for OTEL collector
	// Able to dpeloy to OTEL collector with requested details
	createCustomResource(ctx, otaasPipeLine)

	// Get unique routing key
	routingKey := generateRoutingKey()
	log.Printf(ctx, "Generated Routing Key: [%v]", routingKey)

	// Update mapping with routing key and OTEL collector cluster IP's

	g.JSON(http.StatusCreated, v1.OTaaSRes{
		RoutingKey: routingKey,
	})
}

func createCustomResource(ctx context.Context, pipeline v1.OTaaSPipeline) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube",
			"config"), "Path to a kubeconfig file")

		log.Printf(ctx, "Kubeconfig details : [%+v]", kubeconfig)

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Printf(ctx, "Error loading kubeconfig: %v\n", err)
			return
		}

	}

	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	panic(err.Error())
	//}

	// Load the Kubernetes configuration from the default location or a given kubeconfig file.
	//kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube",
	//	"config"), "Path to a kubeconfig file")

	//log.Infof(ctx, "Kubeconfig details : [%+v]", kubeconfig)

	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//if err != nil {
	//	log.Errorf(ctx, "Error loading kubeconfig: %v\n", err)
	//	return
	//}
	log.Printf(ctx, "config details : [%+v]", config)

	// Create a dynamic client for working with custom resources.
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf(ctx, "Error creating dynamic client: %v\n", err)
		return
	}

	// Define the GVR (GroupVersionResource) for your custom resource.
	gvr := schema.GroupVersionResource{
		Group:    "opentelemetry.io",
		Version:  "v1alpha1",
		Resource: "opentelemetrycollectors",
	}

	// Create an instance of the custom resource.
	customResource := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "opentelemetry.io/v1alpha1",
			"kind":       "OpenTelemetryCollector",
			"metadata": map[string]interface{}{
				"name": pipeline.Name + "-otel-coll-pipeline",
			},
			"spec": map[string]interface{}{
				"config": `receivers:
  otlp:
    protocols:
      grpc:
        endpoint: ${MY_POD_IP}:4317
      http:
        endpoint: ${MY_POD_IP}:4318
  hostmetrics:
    collection_interval: 60s
    scrapers:
      cpu:
      load:
      memory:
      disk:
      filesystem:
      network:
      paging:
      processes:
processors:
  batch:
    send_batch_max_size: 1000
    send_batch_size: 100
    timeout: 10s
  memory_limiter:
    check_interval: 5s
    limit_mib: 2000
exporters:
  logging:
    verbosity: detailed
    sampling_initial: 5
    sampling_thereafter: 200
  prometheusremotewrite:
    endpoint: https://listener-wa.logz.io:8053
    external_labels:
      p8s_logzio_name:` + ` "` + pipeline.Name + `-otaas-metrics"
    headers:
      Authorization: "Bearer gnrnmnGbTQkFljPQBwsLtGOWuMtZDTSl"
service:
  pipelines:
    metrics:
      receivers: [hostmetrics, otlp]
      processors: [batch]
      exporters: [prometheusremotewrite]
`,
			},
		},
	}

	// Apply the custom resource to the cluster.
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err = client.Resource(gvr).Namespace("default").Create(context.TODO(), customResource, metav1.CreateOptions{})
		return err
	})

	if err != nil {
		log.Printf(ctx, "Error applying custom resource: %v\n", err)
		return
	}

	log.Printf(ctx, "Custom resource applied successfully.")
}
func generateRoutingKey() string {
	// Get current timestamp
	timestamp := time.Now().UnixNano()

	// Generate a random number
	randomNumber := rand.Intn(100000)

	// Combine timestamp and random number to create a unique string
	uniqueString := fmt.Sprintf("%d-%d", timestamp, randomNumber)

	// Hash the unique string using SHA-1
	hasher := sha1.New()
	hasher.Write([]byte(uniqueString))
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Return the unique routing key
	return hash
}
