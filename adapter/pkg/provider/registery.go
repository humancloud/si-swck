// Copyright (c) 2019 Stackinsights to present
// All rights reserved
package provider

import (
	"encoding/json"

	"github.com/humancloud/si-cli/pkg/graphql/metrics"
	"k8s.io/apimachinery/pkg/util/wait"
	apiprovider "sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

func (p *externalMetricsProvider) ListAllExternalMetrics() (externalMetricsInfo []apiprovider.ExternalMetricInfo) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for _, md := range p.metricDefines {
		info := apiprovider.ExternalMetricInfo{
			Metric: p.getMetricNameWithNamespace(md.Name),
		}
		externalMetricsInfo = append(externalMetricsInfo, info)
	}
	return
}

func (p *externalMetricsProvider) sync() {
	go wait.Until(func() {
		if err := p.updateMetrics(); err != nil {
			klog.Errorf("failed to update metrics: %v", err)
		}
	}, p.refreshRegistryInterval, wait.NeverStop)
}

func (p *externalMetricsProvider) updateMetrics() error {
	mdd, err := metrics.ListMetrics(p.ctx, p.regex)
	if err != nil {
		return err
	}
	klog.Infof("Get service metrics: %s", display(mdd))
	if len(mdd) > 0 {
		p.lock.Lock()
		defer p.lock.Unlock()
		p.metricDefines = mdd
	}
	return nil
}

func display(data interface{}) string {
	bytes, e := json.Marshal(data)
	if e != nil {
		return "Error"
	}
	return string(bytes)

}
