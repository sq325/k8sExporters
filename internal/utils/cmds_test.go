package utils

import "testing"

func TestJsonStrToMap(t *testing.T) {
	jsonStr := `{"controller-revision-hash":"58bf5dfbd7","k8s-app":"kube-proxy","pod-template-generation":"1"}`
	m, _ := JsonStrToMap(jsonStr)
	for k, v := range m {
		t.Log(k, v)
	}
}

func TestMapToStr(t *testing.T) {
	m := map[string]string{
		"a": "b",
		"c": "d",
	}
	s, _ := MapToStr(m)
	t.Log(s)
}
