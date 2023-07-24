# k8sExporters

## svcPod exporter metrics expamle

```
# HELP kube_pod_labels Contain all labels of pods.
# TYPE kube_pod_labels gauge
kube_pod_labels{labels="{\"addonmanager.kubernetes.io/mode\":\"Reconcile\",\"integration-test\":\"storage-provisioner\"}",name="storage-provisioner",namespace="kube-system"} 1
kube_pod_labels{labels="{\"app\":\"kindnet\",\"controller-revision-hash\":\"78f985b4\",\"k8s-app\":\"kindnet\",\"pod-template-generation\":\"1\",\"tier\":\"node\"}",name="kindnet-lqzpd",namespace="kube-system"} 1
kube_pod_labels{labels="{\"app\":\"prometheus\",\"pod-template-hash\":\"6856764f47\"}",name="prometheus-6856764f47-mw6gc",namespace="monitoring"} 1
kube_pod_labels{labels="{\"component\":\"etcd\",\"tier\":\"control-plane\"}",name="etcd-minikube",namespace="kube-system"} 1
kube_pod_labels{labels="{\"component\":\"kube-apiserver\",\"tier\":\"control-plane\"}",name="kube-apiserver-minikube",namespace="kube-system"} 1
kube_pod_labels{labels="{\"component\":\"kube-controller-manager\",\"tier\":\"control-plane\"}",name="kube-controller-manager-minikube",namespace="kube-system"} 1
kube_pod_labels{labels="{\"component\":\"kube-scheduler\",\"tier\":\"control-plane\"}",name="kube-scheduler-minikube",namespace="kube-system"} 1
kube_pod_labels{labels="{\"controller-revision-hash\":\"58bf5dfbd7\",\"k8s-app\":\"kube-proxy\",\"pod-template-generation\":\"1\"}",name="kube-proxy-49wpd",namespace="kube-system"} 1
kube_pod_labels{labels="{\"k8s-app\":\"kube-dns\",\"pod-template-hash\":\"6d4b75cb6d\"}",name="coredns-6d4b75cb6d-bz268",namespace="kube-system"} 1
kube_pod_labels{labels="{\"k8s-app\":\"metrics-server\",\"pod-template-hash\":\"8595bd7d4c\"}",name="metrics-server-8595bd7d4c-66v84",namespace="kube-system"} 1
kube_pod_labels{labels="{\"run\":\"busybox\"}",name="busybox",namespace="default"} 1
# HELP kube_service_selector kube_service_selector contain all selector of services
# TYPE kube_service_selector gauge
kube_service_selector{name="kube-dns",namespace="kube-system",selector="{\"k8s-app\":\"kube-dns\"}"} 1
kube_service_selector{name="kubernetes",namespace="default",selector=""} 1
kube_service_selector{name="metrics-server",namespace="kube-system",selector="{\"k8s-app\":\"metrics-server\"}"} 1
kube_service_selector{name="prometheus",namespace="monitoring",selector="{\"app\":\"prometheus\"}"} 1
```

---


```mermaid

classDiagram
  class EmptyDir {
    - path   *path.Path
	  - volume *Volume
  }
  class IPodEmptydirs {
    <<interface>>
    + PodName() string
    + PodNamespace() string
    + PodUID() string
    + PodHostIP() string
    + EmptydirListSizeBytes() map[string]int64
    + EmptydirListSizeLimitBytes() map[string]int64
  }
  class PodEmptyDirs {
	  + Pod          *Pod
	  + EmptydirList []*EmptyDir
  }
  class Pod {
    - pod           *coreV1.Pod
  }

  class Volume  {
    - volume *coreV1.Volume
    + Name() string
    + Type() string
  }
  class coreV1Pod {
    + TypeMeta 
	  + Spec PodSpec 
  	+ Status PodStatus 
  }

  class coreV1Volume {
    + Name string
    + VolumeSource
  }
  class Path {
    + AbsPath  string
    + IsDir    bool
    - fileSize int64
  }
  
  EmptyDir o--> Path
  EmptyDir o--> Volume
  coreV1Pod o--> coreV1Volume
  Volume o--> coreV1Volume
  Pod o--> coreV1Pod
  PodEmptyDirs o--> Pod: 组合
  IPodEmptydirs <|.. PodEmptyDirs: implements
  PodEmptyDirs o--> EmptyDir: 组合
```