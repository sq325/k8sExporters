package svcpodCollector

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sq325/k8sExporters/internal/utils"
	"github.com/sq325/k8sExporters/pkg/resource"
)

var (
	SvcMetricName      = "kube_service_selector"
	SvcMetricHelp      = "kube_service_selector contains all selectors services"
	SvcMetricLabelKeys = []string{"namespace", "service", "selector"}
)

var (
	svcCounterVec *prometheus.CounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: SvcMetricName,
		Help: SvcMetricHelp,
	}, SvcMetricLabelKeys)
)

type SvcCollector struct {
	cv     *prometheus.CounterVec
	factor resource.Factor
}

func (c *SvcCollector) Describe(ch chan<- *prometheus.Desc) {
	c.cv.MetricVec.Describe(ch)
}

func (c *SvcCollector) Collect(ch chan<- prometheus.Metric) {
	svcs, err := c.factor.GetResources()
	if err != nil {
		log.Println(err)
		return
	}
	for _, s := range svcs {
		selectorStr, err := utils.MapToStr(s.Selector())
		if err != nil {
			log.Fatal(err)
			return
		}
		svcName := s.Name()
		namespace := s.Namespace()
		values := []string{namespace, svcName, selectorStr}
		c.cv.WithLabelValues(values...).Inc()
	}
	c.cv.Collect(ch)
}

func NewSvcCollector(factor resource.Factor) prometheus.Collector {
	return &SvcCollector{
		cv:     svcCounterVec,
		factor: factor,
	}
}
