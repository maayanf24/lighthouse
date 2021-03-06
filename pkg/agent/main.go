package main

import (
	"flag"

	"github.com/kelseyhightower/envconfig"
	"github.com/submariner-io/lighthouse/pkg/agent/controller"
	lighthousev1 "github.com/submariner-io/lighthouse/pkg/apis/lighthouse.submariner.io/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var (
	masterURL  string
	kubeConfig string
)

func main() {
	agentSpec := controller.AgentSpecification{}
	klog.InitFlags(nil)
	flag.Parse()

	err := envconfig.Process("submariner", &agentSpec)
	if err != nil {
		klog.Fatal(err)
	}
	klog.Infof("AgentSpec: %v", agentSpec)
	if err != nil {
		klog.Fatal(err)
	}
	err = lighthousev1.AddToScheme(scheme.Scheme)
	if err != nil {
		klog.Exitf("Error adding lighthouse V1 to the scheme: %v", err)
	}
	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	klog.Infof("Starting submariner-lighthouse-agent %v", agentSpec)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	lightHouseAgent, err := controller.New(&agentSpec, cfg)

	if err != nil {
		klog.Fatalf("Failed to create lighthouse agent: %v", err)
	}

	if err := lightHouseAgent.Run(stopCh); err != nil {
		klog.Fatalf("Failed to start lighthouse agent: %v", err)
	}

	klog.Info("All controllers stopped or exited. Stopping main loop")
}

func init() {
	flag.StringVar(&kubeConfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
