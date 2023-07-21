package collector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/k8sExporters/internal/utils"
	"github.com/sq325/k8sExporters/pkg/resource"
)

// kube_pod_label metric meta info
var (
	PodMetricName      = "kube_pod_label"
	PodMetricHelp      = "kube_pod_label contains all pods labels"
	PodMetricLabelKeys = []string{"namespace", "pod", "labels"}
)

var (
	podCounterVec *prometheus.CounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: PodMetricName,
		Help: PodMetricHelp,
	}, PodMetricLabelKeys)
)

type PodCollector struct {
	cv     *prometheus.CounterVec
	factor *resource.PodFactor
}

func (c *PodCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cv.MetricVec.Describe(ch)
}

func (c *PodCollector) Collect(ch chan<- prometheus.Metric) {
	pods, err := c.factor.GetPods()
	if err != nil {
		log.Println(err)
		return
	}
	for _, p := range pods {
		labelsStr, err := utils.MapToStr(p.Labels())
		if err != nil {
			log.Fatal(err)
			return
		}
		podName := p.Name()
		namespace := p.Namespace()
		values := []string{namespace, podName, labelsStr}
		c.cv.WithLabelValues(values...).Inc()
	}
	c.cv.Collect(ch)
}

func NewPodCollector(factor *resource.PodFactor) prometheus.Collector {
	return &PodCollector{
		cv:     podCounterVec,
		factor: factor,
	}
}
