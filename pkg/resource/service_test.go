package resource

import (
	"os"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func BenchmarkGetRsources(b *testing.B) {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	// Create a Kubernetes clientset
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, _ := kubernetes.NewForConfig(config)
	factor := NewPodFactor(clientset)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factor.GetResources()
	}
}
