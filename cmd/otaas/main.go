package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"

	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	v1 "github.com/harish-nair-rajagopal/telaas/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
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

	// DELETE pipeline

	// Status Check for livez readyz endpoint
	router.GET("/status", statusCheck)

	router.Run(*listen)
}

func statusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}

func UpdateConfigMap(obj map[string]interface{}, routingKey string, clusterIP string) map[string]interface{} {

	// Update processors in routing table
	type NewVal map[string]interface{}
	exporterVal := fmt.Sprintf("otlp/%s", routingKey)
	AppB := NewVal{"value": routingKey, "exporters": []string{exporterVal}}
	m := obj["processors"].(map[interface{}]interface{})
	for k := range m {
		fmt.Println(k)
		if k == "routing" {
			tab := m[k].(map[interface{}]interface{})
			for t := range tab {
				if t == "table" {
					fmt.Println(tab[t])
					tab[t] = append(tab[t].([]interface{}), AppB)
				}
			}
		}
	}

	// Update Exporters to the specific endoint
	type exp map[string]interface{}
	exp1 := exp{
		"endpoint": fmt.Sprintf("%s:4317", clusterIP),
		"tls": map[string]interface{}{
			"insecure": true,
		},
	}
	z := obj["exporters"].(map[interface{}]interface{})
	for i := range z {
		fmt.Println(i, z[i])
		z[exporterVal] = exp1

	}

	//Update the service pipelines

	n := obj["service"].(map[interface{}]interface{})
	for j := range n {
		pipe := n[j].(map[interface{}]interface{})["metrics"].(map[interface{}]interface{})
		for y := range pipe {
			if y == "exporters" {
				pipe[y] = append(pipe[y].([]interface{}), fmt.Sprintf("otlp/%s", routingKey))
				break
			}
		}
	}

	return obj

}

func UpdateListener(ctx context.Context, pipeline v1.OTaaSPipeline, routingKey string, conf *rest.Config) {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube",
	// 		"config"), "Path to a kubeconfig file")

	// 	log.Printf("Kubeconfig details : [%+v]", kubeconfig)

	// 	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// 	if err != nil {
	// 		log.Printf("Error loading kubeconfig: %v\n", err)
	// 		return
	// 	}

	// }

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic(err.Error())
	}

	// Create a object to hold old configMap data
	obj := make(map[string]interface{})

	ListenerOtelConfName := "collector-config"

	oldconfigMap, err := clientset.CoreV1().ConfigMaps("default").Get(context.TODO(), ListenerOtelConfName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("configMap not found!!")
	} else {
		data := oldconfigMap.Data["collector.yaml"]
		fmt.Println("check", data)
		err := yaml.Unmarshal([]byte(data), obj)
		if err != nil {
			fmt.Print("err")
		}
		// Getting Cluster IP of the service
		serviceName := fmt.Sprintf("%s-otel-coll-pipeline-collector", pipeline.Name)
		// fmt.Println("-----------", serviceName)
		// time.Sleep(20 * time.Second)
		// svcList, err := clientset.CoreV1().Services("default").Get(context.TODO(), serviceName, metav1.GetOptions{})
		// fmt.Println(svcList)
		// clusterIP := svcList.Spec.ClusterIP
		// fmt.Println(clusterIP)

		obj = UpdateConfigMap(obj, routingKey, serviceName)
		res, err := yaml.Marshal(obj)
		fmt.Println(string(res))
		// configMapData := make(map[string]string, 0)
		// configMapData["collector.yaml"] = string(res)
		oldconfigMap.Data["collector.yaml"] = string(res)

		// Updating the config Map
		clientset.CoreV1().ConfigMaps("default").Update(context.TODO(), oldconfigMap, metav1.UpdateOptions{})

	}
}

func CreateOTELPipeline(g *gin.Context) {
	ctx := g.Request.Context()

	var otaasPipeLine v1.OTaaSPipeline

	if err := g.BindJSON(&otaasPipeLine); err != nil {
		log.Printf("error: %v", err)
		g.JSON(http.StatusBadRequest, err.Error())

		return
	}

	log.Printf("creating OtelPipe with details [%+v]", otaasPipeLine)

	//Construct Custom Resource for OTEL collector
	// Able to dpeloy to OTEL collector with requested details
	conf := createCustomResource(ctx, otaasPipeLine)

	// Get unique routing key
	routingKey := generateRoutingKey()
	log.Printf("Generated Routing Key: [%v]", routingKey)

	// Update mapping with routing key and OTEL collector cluster IP's
	UpdateListener(ctx, otaasPipeLine, routingKey, conf)

	g.JSON(http.StatusCreated, v1.OTaaSRes{
		RoutingKey: routingKey,
	})
}

func createCustomResource(ctx context.Context, pipeline v1.OTaaSPipeline) *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube",
			"config"), "Path to a kubeconfig file")

		log.Printf("Kubeconfig details : [%+v]", kubeconfig)

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Printf("Error loading kubeconfig: %v\n", err)
			panic(err.Error())
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
	log.Printf("config details : [%+v]", config)

	// Create a dynamic client for working with custom resources.
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Error creating dynamic client: %v\n", err)
		panic(err.Error())
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
      receivers: [otlp]
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
		log.Printf("Error applying custom resource: %v\n", err)
		panic(err.Error())
	}

	log.Printf("Custom resource applied successfully.")
	return config

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
