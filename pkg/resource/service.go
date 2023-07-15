package resource

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Service implement Resource interface
// Service define a service resource
type Svc struct {
	name, namespace string
	kind            string // title style
	selector        map[string]string
	labels          map[string]string
}
type Svcs []*Svc

func NewSvc(name, namespaces string, selector map[string]string, labels map[string]string) *Svc {
	return &Svc{
		name:      name,
		namespace: namespaces,
		kind:      "Service",
		selector:  selector,
		labels:    labels,
	}
}

func (s *Svc) Name() string {
	return s.name
}
func (s *Svc) Namespace() string {
	return s.namespace
}

func (s *Svc) Kind() string {
	return s.kind
}

func (s *Svc) Selector() map[string]string {
	return s.selector
}

func (s *Svc) Labels() map[string]string {
	return s.labels
}

// SvcFactor implement Factor interface
type SvcFactor struct {
	ClientSet *kubernetes.Clientset
}

func NewSvcFactor(clientSet *kubernetes.Clientset) Factor {
	return &SvcFactor{
		ClientSet: clientSet,
	}
}

func (s *SvcFactor) GetResources() (Resources, error) {
	svcList, err := s.ClientSet.CoreV1().Services("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	var svcs Resources
	for _, nsvc := range svcList.Items {
		svc := NewSvc(nsvc.Name, nsvc.Namespace, nsvc.Spec.Selector, nsvc.Labels)
		svcs = append(svcs, svc)
	}

	if len(svcs) == 0 {
		log.Println("No services found in cluster")
		return nil, nil
	}
	return svcs, nil
}
