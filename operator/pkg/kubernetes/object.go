// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package kubernetes

import (
	"encoding/hex"
	"encoding/json"
	"hash/fnv"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func getVersion(o *unstructured.Unstructured, key string) string {
	ann := o.GetAnnotations()
	if ann == nil {
		return ""
	}
	return ann[key]
}

func setVersion(o *unstructured.Unstructured, key string, version string) {
	ann := make(map[string]string)
	for k, v := range o.GetAnnotations() {
		ann[k] = v
	}
	ann[key] = version
	o.SetAnnotations(ann)
}

func hash(o *unstructured.Unstructured) (string, error) {
	forHashing := make(map[string]interface{})
	for field, contents := range o.Object {
		if !isMeta(field) {
			forHashing[field] = contents
		}
	}
	if len(forHashing) == 0 {
		forHashing = map[string]interface{}{"no-contents-use-meta": metaHash(o)}
	}
	bytes, err := json.Marshal(forHashing)
	if err != nil {
		return "", err
	}
	return alphanumeric(bytes), nil
}

func metaHash(o *unstructured.Unstructured) string {
	if o == nil {
		return "-"
	}
	kind := o.GetObjectKind().GroupVersionKind().Kind
	namespace := o.GetNamespace()
	switch kind {
	case "ClusterRole", "ClusterRoleBinding":
		namespace = ""
	}
	return strings.Join([]string{kind, namespace, o.GetName()}, ":")
}

func isMeta(name string) bool {
	switch name {
	case "kind", "apiVersion", "metadata":
		return true

	default:
		return false
	}
}

func alphanumeric(in []byte) string {
	if in == nil {
		return ""
	}

	hash := fnv.New64a()
	_, _ = hash.Write(in)
	out := hash.Sum(make([]byte, 0, 8))
	return hex.EncodeToString(out)
}
