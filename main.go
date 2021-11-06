package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL  string
	kubeconfig string
	addr       = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {

	kubeCfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	log.Info("Built config from flags...")

	kubeClient, err := kubernetes.NewForConfig(kubeCfg)
	if err != nil {
		log.Fatalf("Error building watcher clientset: %s", err.Error())
	}

	// Prometheus Http handler
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
	// Prometheus metrics
	cns_pod_current_count := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cns_current_pod_count",
		Help: "Counts the number of pods running in the cluster",
	})
	// Runtime loop

	for {
		// Fetch pods in cluster
		podList, err := kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Errorf("Error fetching pods: %s", err.Error())
		} else {
			cns_pod_current_count.Set(float64(len(podList.Items)))
		}
		// Update metric

		time.Sleep(time.Second * 10)
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}
