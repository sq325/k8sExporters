package resource

// Resource is the interface a k8s resource must implement
type Resource interface {
	Name() string
	Namespace() string
	Kind() string
	Labels() map[string]string
	Selector() map[string]string
}

type Resources []Resource

// Factor is the interface a resource factor must implement
type Factor interface {
	GetResources() (Resources, error)
}
